package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"time"
	"user-mananger/config"
)

func Connect(ctx context.Context, cfg config.KafkaConfig) (*kafka.Conn, error) {
	for {
		conn, err := kafka.DialLeader(ctx, "tcp", cfg.Broker, cfg.Topic, 0)
		if err != nil {
			slog.Warn("connected to kafka failed, retry", slog.String("err", err.Error()))
			time.Sleep(1 * time.Second)
		} else {
			slog.Info("connected to kafka")
			return conn, nil
		}
	}
}
