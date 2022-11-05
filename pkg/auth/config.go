package auth

import (
	"flag"
	"fmt"
)

// Config for a gateway
type Config struct {
	JwtSecret     string
	ExtraHeaders  string
	TenantName    string
	TenantIDClaim string

	JwksURL             string
	JwksRefreshEnabled  bool
	JwksRefreshInterval int
	JwksRefreshTimeout  int
}

// RegisterFlags adds the flags required to config this package's Config struct
func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&cfg.TenantName, "gateway.auth.tenant-name", "", "Tenant name to use when jwt auth disabled")
	f.StringVar(&cfg.JwtSecret, "gateway.auth.jwt-secret", "", "Secret to sign JSON Web Tokens")
	f.StringVar(&cfg.ExtraHeaders, "gateway.auth.jwt-extra-headers", "", "A comma separated list of additional headers to scan for JSON Web Tokens presence")
	f.StringVar(&cfg.TenantIDClaim, "gateway.auth.tenant-id-claim", "tenant_id", "The name of the Tenant ID Claim. Defaults to tenant_id")

	f.StringVar(&cfg.JwksURL, "gateway.auth.jwks-url", "", "The URL to load the JWKS (JSON Web Key Set) from")
	f.BoolVar(&cfg.JwksRefreshEnabled, "gateway.auth.jwks-refresh-enabled", false, "Enable the JWKS background refresh. (Default: false)")
	f.IntVar(&cfg.JwksRefreshInterval, "gateway.auth.jwks-refresh-interval", 60, "The JWKS background refresh interval in minutes. (Defaults: 60 minutes)")
	f.IntVar(&cfg.JwksRefreshTimeout, "gateway.auth.jwks-refresh-timeout", 30, "The JWKS background refresh timeout in seconds. (Defaults: 30 seconds)")
}

// Validate given config parameters. Returns nil if everything is fine
func (cfg *Config) Validate() error {
	if cfg.JwtSecret == "" && cfg.JwksURL == "" && cfg.TenantName == "" {
		return fmt.Errorf("you must set -gateway.auth.jwt-secret and/or -gateway.auth.jwks-url or -gateway.auth.tenantName")
	}

	if cfg.JwksRefreshInterval <= 0 {
		return fmt.Errorf("JWKS background refresh interval must positive. Given is '%v'", cfg.JwksRefreshInterval)
	}

	if cfg.JwksRefreshTimeout <= 0 {
		return fmt.Errorf("JWKS background refresh timeout must positive. Given is '%v'", cfg.JwksRefreshTimeout)
	}

	return nil
}
