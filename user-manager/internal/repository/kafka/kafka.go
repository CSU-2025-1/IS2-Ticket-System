package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"user-mananger/config"
)

func Connect(ctx context.Context, cfg config.KafkaConfig) (*kafka.Conn, error) {
	return kafka.DialLeader(ctx, "tcp", cfg.Address, cfg.Topic, 0)
}
