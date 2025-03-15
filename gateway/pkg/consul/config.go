package consul

import "time"

// Config is a configuration parameters model for Consul client
type Config struct {
	ConsulAddress       string
	HealthCheckInterval time.Duration
	HealthCheckTimeout  time.Duration
}
