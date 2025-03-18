package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Config is a general configuration file for all systems of service
type Config struct {
	Redis    Redis    `env-prefix:"REDIS_"`
	Registry Registry `env-prefix:"REGISTRY_"`
	Consul   Consul   `env-prefix:"CONSUL_"`
	Proxy    Proxy    `env-prefix:"PROXY_"`
}

// LoadConfig load config form envs
func LoadConfig() (*Config, error) {
	cfg := new(Config)

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("Config.LoadConfig: %w", err)
	}

	return cfg, nil
}

// LoadDotEnv load config from .env
func LoadDotEnv() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("Config.LoadDotEnv: %w", err)
	}

	return nil
}
