package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	MaxConns int32  `yaml:"max_conns"`
}

func (c *Config) ToConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}

func Connect(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	connStr := cfg.ToConnString()
	pgxConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse conn str %s: %w", connStr, err)
	}

	pgxConfig.MaxConns = cfg.MaxConns

	return pgxpool.NewWithConfig(ctx, pgxConfig)
}
