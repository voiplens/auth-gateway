package app

import (
	"net/http"

	"github.com/celest-io/mimir-gateway/pkg/auth"
	"github.com/celest-io/mimir-gateway/pkg/proxy"

	"github.com/cortexproject/cortex/pkg/util/log"
	klog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/weaveworks/common/server"
)

type Gateway struct {
	cfg                Config
	authCfg            auth.Config
	distributorProxy   *proxy.Proxy
	queryFrontendProxy *proxy.Proxy
	rulerProxy         *proxy.Proxy
	alertManagerProxy  *proxy.Proxy
	server             *server.Server
}

// New instantiates a new Gateway
func New(cfg Config, authCfg auth.Config, svr *server.Server) (*Gateway, error) {
	// Initialize reverse proxy for each upstream target service
	distributor, err := proxy.NewProxy(cfg.DistributorAddress, "distributor")
	if err != nil {
		return nil, err
	}
	queryFrontend, err := proxy.NewProxy(cfg.QueryFrontendAddress, "query-frontend")
	if err != nil {
		return nil, err
	}
	ruler, err := proxy.NewProxy(cfg.RulerAddress, "ruler")
	if err != nil {
		return nil, err
	}
	alertManager, err := proxy.NewProxy(cfg.AlertManagerAddress, "ruler")
	if err != nil {
		return nil, err
	}

	return &Gateway{
		cfg:                cfg,
		authCfg:            authCfg,
		distributorProxy:   distributor,
		queryFrontendProxy: queryFrontend,
		rulerProxy:         ruler,
		alertManagerProxy:  alertManager,
		server:             svr,
	}, nil
}

// Start initializes the Gateway and starts it
func (g *Gateway) Start() {
	g.registerRoutes()
}

// RegisterRoutes binds all to be piped routes to their handlers
func (g *Gateway) registerRoutes() {
	authenticateTenant := auth.NewAuthenticationMiddleware(g.authCfg)

	g.server.HTTP.Path("/all_user_stats").HandlerFunc(g.distributorProxy.Handler)
	g.server.HTTP.Path("/api/prom/push").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.distributorProxy.Handler)))
	g.server.HTTP.Path("/api/v1/push").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.distributorProxy.Handler)))

	g.server.HTTP.PathPrefix("/prometheus/api/v1/alerts").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.rulerProxy.Handler)))
	g.server.HTTP.PathPrefix("/prometheus/api/v1/rules").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.rulerProxy.Handler)))
	g.server.HTTP.PathPrefix("/prometheus/config/v1/rules").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.rulerProxy.Handler)))
	g.server.HTTP.PathPrefix("/api/prom/api/v1/alerts").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.rulerProxy.Handler)))
	g.server.HTTP.PathPrefix("/api/prom/api/v1/rules").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.rulerProxy.Handler)))
	g.server.HTTP.PathPrefix("/api/v1/rules").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.rulerProxy.Handler)))
	g.server.HTTP.PathPrefix("/api/prom/rules").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.rulerProxy.Handler)))

	g.server.HTTP.PathPrefix("/prometheus/api").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.queryFrontendProxy.Handler)))
	g.server.HTTP.PathPrefix("/api/v1/user_stats").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.queryFrontendProxy.Handler)))

	g.server.HTTP.PathPrefix("/api/v1/alerts").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.alertManagerProxy.Handler)))
	g.server.HTTP.PathPrefix("/alertmanager").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.alertManagerProxy.Handler)))
	g.server.HTTP.PathPrefix("/multitenant_alertmanager/status").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.alertManagerProxy.Handler)))
	g.server.HTTP.PathPrefix("/api/v1/alerts").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.alertManagerProxy.Handler)))
	g.server.HTTP.PathPrefix("/api/prom/alertmanager").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.alertManagerProxy.Handler)))

	g.server.HTTP.PathPrefix("/api").Handler(authenticateTenant.Wrap(http.HandlerFunc(g.queryFrontendProxy.Handler)))
	g.server.HTTP.Path("/health").HandlerFunc(g.healthCheck)
	g.server.HTTP.PathPrefix("/").HandlerFunc(g.notFoundHandler)
}

func (g *Gateway) healthCheck(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "ok", http.StatusOK)
}

func (g *Gateway) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	logger := klog.With(log.WithContext(r.Context(), log.Logger), "ip_address", r.RemoteAddr)
	level.Info(logger).Log("msg", "no request handler defined for this route", "route", r.RequestURI)
	http.NotFound(w, r)
}
