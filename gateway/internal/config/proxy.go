package config

// Proxy is a configuration parameters model for proxy server
type Proxy struct {
	BalancerAlgorithm string `env:"BALANCER_ALGORITHM" envDefault:"round_robin"`
	EnableCaching     bool   `env:"ENABLE_CACHING" envDefault:"false"`
	LaunchedPort      uint16 `env:"LAUNCHED_PORT" envDefault:"80"`
}
