package config

type Http struct {
	Port    uint16 `env:"PORT"`
	Address string `env:"ADDRESS"`
}
