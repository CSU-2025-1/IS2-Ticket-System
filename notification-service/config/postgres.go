package config

type Postgres struct {
	ConnectionString string `env:"CONNECTION_STRING"`
}
