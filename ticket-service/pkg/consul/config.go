package consul

import "time"

// Config is a configuration parameters model for Consul client
type Config struct {
	ConsulAddress       string        `yaml:"consul_address"`
	HealthCheckInterval time.Duration `yaml:"health_check_interval"`
	HealthCheckTimeout  time.Duration `yaml:"health_check_timeout"`
}
