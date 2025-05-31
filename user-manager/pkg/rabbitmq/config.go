package rabbitmq

import "fmt"

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Queue    string `yaml:"queue"`
	PoolSize int    `yaml:"pool_size"`
}

func (c *Config) ToDSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d", c.User, c.Password, c.Host, c.Port)
}
