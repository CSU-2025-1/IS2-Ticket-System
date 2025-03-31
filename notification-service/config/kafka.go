package config

type Kafka struct {
	Address             string `env:"ADDRESS"`
	GroupID             string `env:"GROUP_ID"`
	AutoOffsetReset     string `env:"AUTO_OFFSET_RESET"`
	TicketCreationTopic string `env:"TICKET_CREATION_TOPIC"`
}
