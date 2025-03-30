package kafka

// ConsumerConfig is configuration params for Kafka
type ConsumerConfig struct {
	Address         string
	GroupID         string
	AutoOffsetReset string
	Topic           string
}
