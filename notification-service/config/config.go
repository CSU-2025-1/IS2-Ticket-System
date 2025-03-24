package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Config is a model for configuration params
type Config struct {
	Mail     Mail     `env-prefix:"MAIL_"`
	Kafka    Kafka    `env-prefix:"KAFKA_"`
	Postgres Postgres `env-prefix:"POSTGRES_"`
	Http     Http     `env-prefix:"HTTP_"`
	Consul   Consul   `env-prefix:"CONSUL_"`
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
