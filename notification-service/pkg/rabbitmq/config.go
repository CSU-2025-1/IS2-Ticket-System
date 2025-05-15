package rabbitmq

import "fmt"

type Config struct {
	Host     string `env:"HOST" env-default:"rabbitmq"`
	Port     int    `env:"PORT" env-default:"5672"`
	User     string `env:"USER" env-default:"admin"`
	Password string `env:"PASSWORD" env-default:"password"`
	Queue    string `env:"QUEUE" env-default:"new_tickets"`
}

func (c *Config) ToDSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d", c.User, c.Password, c.Host, c.Port)
}
