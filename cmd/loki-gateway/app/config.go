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
	QuerierAddress       string
	RulerAddress         string
}

// RegisterFlags adds the flags required to config this package's Config struct
func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&cfg.DistributorAddress, "gateway.distributor.address", "", "Upstream HTTP URL for Loki Distributor")
	f.StringVar(&cfg.QueryFrontendAddress, "gateway.query-frontend.address", "", "Upstream HTTP URL for Loki Query Frontend")
	f.StringVar(&cfg.RulerAddress, "gateway.ruler.address", "", "Upstream HTTP URL for Loki Ruler")
	f.StringVar(&cfg.QuerierAddress, "gateway.querier.address", "", "Upstream HTTP URL for Loki Querier")
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

	if cfg.QuerierAddress == "" {
		return fmt.Errorf("you must set -gateway.querier.address")
	}

	if !strings.HasPrefix(cfg.QuerierAddress, "http") {
		return fmt.Errorf("querier address must start with a valid scheme (http/https). Given is '%v'", cfg.QuerierAddress)
	}

	return nil
}
