package main

import (
	envstruct "code.cloudfoundry.org/go-envstruct"
	sharedtls "github.com/cloudfoundry/metric-store-release/src/pkg/tls"
)

// Config is the configuration for a MetricStore Gateway.
type Config struct {
	Addr            string `env:"ADDR, required, report"`
	MetricStoreAddr string `env:"METRIC_STORE_ADDR, required, report"`
	HealthPort      int    `env:"HEALTH_PORT, report"`
	ProxyCertPath   string `env:"PROXY_CERT_PATH, report"`
	ProxyKeyPath    string `env:"PROXY_KEY_PATH, report"`
	TLS             sharedtls.TLS
}

// LoadConfig creates Config object from environment variables
func LoadConfig() (*Config, error) {
	c := Config{
		Addr:            "localhost:8081",
		HealthPort:      6063,
		MetricStoreAddr: "localhost:8080",
	}

	if err := envstruct.Load(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
