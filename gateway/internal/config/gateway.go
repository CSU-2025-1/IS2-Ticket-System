package config

import "time"

type Gateway struct {
	HealthCheckInterval time.Duration
	Port                int
	RouteProxyMethod    string
}
