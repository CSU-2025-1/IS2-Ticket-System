package config

import "time"

// Consul is a configuration parameters model for consul registry
type Consul struct {
	ConsulAddress       string        `env:"ADDRESS"`
	HealthCheckInterval time.Duration `env:"HEALTH_CHECK_INTERVAL"`
	HealthCheckTimeout  time.Duration `env:"HEALTH_CHECK_TIMEOUT"`
}
