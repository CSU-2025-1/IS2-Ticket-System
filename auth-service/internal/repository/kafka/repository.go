package kafka

import (
	"auth-service/config"
	"auth-service/internal/domain/entity"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"time"
)

type KafkaUserStream struct {
	userChan chan *entity.User
	reader   *kafka.Reader
}

func NewKafkaUserStream(cfg *config.KafkaConfig) *KafkaUserStream {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.Broker},
		Topic:    cfg.Topic,
		GroupID:  cfg.GroupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second, // Таймаут подключения
			DualStack: true,
		},
	})

	return &KafkaUserStream{
		reader: reader,
	}
}

func (k *KafkaUserStream) GetNewUserStream(ctx context.Context) <-chan *entity.User {
	if k.userChan != nil {
		return k.userChan
	}

	k.userChan = make(chan *entity.User)

	go func() {
		defer close(k.userChan)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := k.reader.ReadMessage(ctx)
				if err != nil {
					slog.Error("Failed to read message from Kafka",
						slog.String("error", err.Error()),
					)
					continue
				}

				var user User
				if err := json.Unmarshal(msg.Value, &user); err != nil {
					slog.Error("Failed to unmarshal Kafka message",
						slog.String("error", err.Error()),
					)
					continue
				}

				en, err := user.ToEntity()
				if err != nil {
					slog.Error("Failed to convert Kafka message to Entity",
						slog.String("error", err.Error()),
					)
				}

				k.userChan <- en
			}
		}
	}()

	return k.userChan
}

func (k *KafkaUserStream) Close() error {
	return k.reader.Close()
}
