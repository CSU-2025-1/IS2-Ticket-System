package config

import "time"

// Registry is a configuration params model for service addresses registry
type Registry struct {
	ActualizingInterval time.Duration `env:"ACTUALIZE_INTERVAL" env-default:"10s"`
}
