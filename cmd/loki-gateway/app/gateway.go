package app

import (
	"net/http"

	"github.com/celest-io/auth-gateway/pkg/auth"
	"github.com/celest-io/auth-gateway/pkg/proxy"
	"github.com/celest-io/auth-gateway/pkg/util/log"

	klog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/weaveworks/common/server"
)

// Gateway hosts a reverse proxy for each upstream cortex service we'd like to tunnel after successful authentication
type Gateway struct {
	authCfg            auth.Config
	distributorProxy   *proxy.Proxy
	queryFrontendProxy *proxy.Proxy
	rulerProxy         *proxy.Proxy
	querierProxy       *proxy.Proxy
	server             *server.Server
	logger             klog.Logger
}

// NewGateway instantiates a new Gateway
func NewGateway(gatewayCfg Config, authCfg auth.Config, svr *server.Server, logger klog.Logger) (*Gateway, error) {
	// Initialize reverse proxy for each upstream target service
	distributor, err := proxy.NewProxy(gatewayCfg.DistributorAddress, "distributor")
	if err != nil {
		return nil, err
	}
	queryFrontend, err := proxy.NewProxy(gatewayCfg.QueryFrontendAddress, "query-frontend")
	if err != nil {
		return nil, err
	}
	ruler, err := proxy.NewProxy(gatewayCfg.RulerAddress, "ruler")
	if err != nil {
		return nil, err
	}
	querier, err := proxy.NewProxy(gatewayCfg.QuerierAddress, "querier")
	if err != nil {
		return nil, err
	}

	return &Gateway{
		authCfg:            authCfg,
		distributorProxy:   distributor,
		queryFrontendProxy: queryFrontend,
		rulerProxy:         ruler,
		querierProxy:       querier,
		server:             svr,
		logger:             logger,
	}, nil
}

// Start initializes the Gateway and starts it
func (g *Gateway) Start() {
	g.registerRoutes()
}

// RegisterRoutes binds all to be piped routes to their handlers
func (g *Gateway) registerRoutes() {
	authenticateTenant := auth.NewAuthenticationMiddleware(g.authCfg, g.logger)

	g.server.HTTP.Path("/all_user_stats").HandlerFunc(g.distributorProxy.Handler)
	g.server.HTTP.Path("/loki/api/v1/push").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.distributorProxy.Handler)))
	g.server.HTTP.Path("/api/prom/push").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.distributorProxy.Handler)))

	g.server.HTTP.Path("/loki/api/v1/tail").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.querierProxy.Handler)))
	g.server.HTTP.Path("/api/prom/tail").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.querierProxy.Handler)))

	g.server.HTTP.PathPrefix("/prometheus/api/v1/alerts").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.rulerProxy.Handler)))
	g.server.HTTP.PathPrefix("/prometheus/api/v1/rules").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.rulerProxy.Handler)))
	g.server.HTTP.PathPrefix("/api/prom/rules").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.rulerProxy.Handler)))
	g.server.HTTP.PathPrefix("/api/prom/alerts").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.rulerProxy.Handler)))

	g.server.HTTP.PathPrefix("/loki/api/").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.queryFrontendProxy.Handler)))
	g.server.HTTP.PathPrefix("/api/prom/").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.queryFrontendProxy.Handler)))
	g.server.HTTP.Path("/health").HandlerFunc(g.healthCheck)
	g.server.HTTP.PathPrefix("/").HandlerFunc(g.notFoundHandler)
}

func (g *Gateway) healthCheck(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "ok", http.StatusOK)
}

func (g *Gateway) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	logger := klog.With(log.WithContext(r.Context(), g.logger), "ip_address", r.RemoteAddr)
	level.Info(logger).Log("msg", "no request handler defined for this route", "route", r.RequestURI)
	http.NotFound(w, r)
}
