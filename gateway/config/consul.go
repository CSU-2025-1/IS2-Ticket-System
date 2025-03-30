package config

import "time"

// Consul is a configuration parameters model for consul registry
type Consul struct {
	ConsulAddress       string        `env:"ADDRESS" envDefault:"consul:8500"`
	HealthCheckInterval time.Duration `env:"HEALTH_CHECK_INTERVAL" envDefault:"10s"`
	HealthCheckTimeout  time.Duration `env:"HEALTH_CHECK_TIMEOUT" envDefault:"30s"`
}
