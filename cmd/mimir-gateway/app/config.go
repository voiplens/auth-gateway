package app

import (
	"flag"
	"fmt"
	"strings"
)

// Config for a gateway
type Config struct {
	DistributorAddress   string
	QueryFrontendAddress string
	RulerAddress         string
	AlertManagerAddress  string
}

// RegisterFlags adds the flags required to config this package's Config struct
func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&cfg.DistributorAddress, "gateway.distributor.address", "", "Upstream HTTP URL for Mimir Distributor")
	f.StringVar(&cfg.QueryFrontendAddress, "gateway.query-frontend.address", "", "Upstream HTTP URL for Mimir Query Frontend")
	f.StringVar(&cfg.RulerAddress, "gateway.ruler.address", "", "Upstream HTTP URL for Mimir Query Frontend")
	f.StringVar(&cfg.AlertManagerAddress, "gateway.alertmanager.address", "", "Upstream HTTP URL for Mimir AlertManager")
}

// Validate given config parameters. Returns nil if everything is fine
func (cfg *Config) Validate() error {
	if cfg.DistributorAddress == "" {
		return fmt.Errorf("you must set -gateway.distributor.address")
	}

	if !strings.HasPrefix(cfg.DistributorAddress, "http") {
		return fmt.Errorf("distributor address must start with a valid scheme (http/https). Given is '%v'", cfg.DistributorAddress)
	}

	if cfg.QueryFrontendAddress == "" {
		return fmt.Errorf("you must set -gateway.query-frontend.address")
	}

	if !strings.HasPrefix(cfg.QueryFrontendAddress, "http") {
		return fmt.Errorf("query frontend address must start with a valid scheme (http/https). Given is '%v'", cfg.QueryFrontendAddress)
	}

	if cfg.RulerAddress == "" {
		return fmt.Errorf("you must set -gateway.ruler.address")
	}

	if !strings.HasPrefix(cfg.RulerAddress, "http") {
		return fmt.Errorf("ruler address must start with a valid scheme (http/https). Given is '%v'", cfg.RulerAddress)
	}

	if cfg.AlertManagerAddress == "" {
		return fmt.Errorf("you must set -gateway.alertmanager.address")
	}

	if !strings.HasPrefix(cfg.AlertManagerAddress, "http") {
		return fmt.Errorf("The AlertManager address must start with a valid scheme (http/https). Given is '%v'", cfg.AlertManagerAddress)
	}

	return nil
}
