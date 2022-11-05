package main

import (
	"flag"
	"net/http"

	"github.com/celest-io/auth-gateway/cmd/mimir-gateway/app"
	"github.com/celest-io/auth-gateway/pkg/auth"

	"github.com/cortexproject/cortex/pkg/util/log"
	"github.com/grafana/dskit/flagext"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/weaveworks/common/middleware"
	"github.com/weaveworks/common/server"
	"github.com/weaveworks/common/tracing"
	"google.golang.org/grpc"
)

func main() {
	operationNameFunc := nethttp.OperationNameFunc(func(r *http.Request) string {
		return r.URL.RequestURI()
	})

	var (
		serverCfg = server.Config{
			MetricsNamespace: "cortex_gateway",
			HTTPMiddleware: []middleware.Interface{
				middleware.Func(func(handler http.Handler) http.Handler {
					return nethttp.Middleware(opentracing.GlobalTracer(), handler, operationNameFunc)
				}),
			},
			GRPCMiddleware: []grpc.UnaryServerInterceptor{
				middleware.ServerUserHeaderInterceptor,
			},
		}
		gatewayCfg app.Config
		authCfg    auth.Config
	)

	flagext.RegisterFlags(&serverCfg, &gatewayCfg, &authCfg)
	flag.Parse()

	log.InitLogger(&serverCfg)

	// Must be done after initializing the logger, otherwise no log message is printed
	err := gatewayCfg.Validate()
	log.CheckFatal("validating gateway config", err)

	err = authCfg.Validate()
	log.CheckFatal("validating authentication config", err)

	// Setting the environment variable JAEGER_AGENT_HOST enables tracing
	trace, err := tracing.NewFromEnv("cortex-gateway")
	log.CheckFatal("initializing tracing", err)
	defer trace.Close()

	svr, err := server.New(serverCfg)
	log.CheckFatal("initializing server", err)
	defer svr.Shutdown()

	// Setup proxy and register routes
	gateway, err := app.NewGateway(gatewayCfg, authCfg, svr)
	log.CheckFatal("initializing gateway", err)
	gateway.Start()

	err = svr.Run()
	if err != nil {
		log.CheckFatal("Error running server gateway", err)
	}
}
