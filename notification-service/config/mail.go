package config

// Mail is a configuration params model for sending emails
type Mail struct {
	Sender   string `env:"FROM" env-default:"info@kforge.ru"`
	Host     string `env:"HOST" env-default:"connect.smtp.bz"`
	Port     int    `env:"PORT" env-default:"2525"`
	Password string `env:"PASSWORD" env-default:"qwerty12345"`
}
