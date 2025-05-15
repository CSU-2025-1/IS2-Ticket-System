package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"user-mananger/internal/repository/postgres"
	"user-mananger/pkg/consul"
	"user-mananger/pkg/rabbitmq"
)

type Config struct {
	Database *postgres.Config `yaml:"database"`
	Server   *HttpConfig      `yaml:"http"`
	Grpc     *GrpcConfig      `yaml:"grpc"`
	Rabbit   *rabbitmq.Config `yaml:"rabbit"`
	Consul   *consul.Config   `yaml:"consul"`
}

type KafkaConfig struct {
	Broker string `yaml:"broker"`
	Topic  string `yaml:"topic"`
}

type HttpConfig struct {
	Address string `yaml:"address"`
}

type GrpcConfig struct {
	Address string `yaml:"address"`
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
