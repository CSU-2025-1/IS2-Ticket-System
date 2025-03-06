package config

import (
	"auth-service/internal/repository/postgres"
	"auth-service/pkg/hydra"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Config struct {
	Database *postgres.Config `yaml:"database"`
	Hydra    *hydra.Config    `yaml:"hydra"`
	Server   *ServerConfig    `yaml:"server"`
	Hash     *HashConfig      `yaml:"hash"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type HashConfig struct {
	Salt string `yaml:"salt"`
	Cost int    `yaml:"cost"`
}

func Parse(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}

	in, err := io.ReadAll(file)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(in, &cfg); err != nil {
		return Config{}, fmt.Errorf("unmurshal config error: %w", err)
	}

	return cfg, nil
}
