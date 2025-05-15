package rabbitmq

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type JsonQueueAdapter[T any] struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewJsonQueueAdapter[T any](cfg Config) (*JsonQueueAdapter[T], error) {
	conn, err := amqp.Dial(cfg.ToDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %s", err)
	}

	return &JsonQueueAdapter[T]{
		conn:    conn,
		channel: ch,
		queue:   cfg.Queue,
	}, nil
}

func (j *JsonQueueAdapter[T]) Consume(opts ...ConsumeOption) (<-chan T, error) {
	opt := DefaultConsumeOption
	if len(opts) > 0 {
		opt = opts[0]
	}

	msgs, err := j.channel.Consume(
		j.queue,
		opt.Consumer,
		opt.AutoAck,
		opt.Exclusive,
		opt.NoLocal,
		opt.NoWait,
		opt.Args,
	)
	if err != nil {
		return nil, err
	}

	channel := make(chan T)
	go func() {
		for d := range msgs {
			var msg T
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				slog.Error("[ Rabbit MQ failed unmarshal json message ]", slog.String("err", err.Error()))
				continue
			}
			channel <- msg
			d.Ack(false)
		}
	}()

	return channel, nil
}

func (j *JsonQueueAdapter[T]) Publish(msg T, opts ...PublishOption) error {
	opt := DefaultPublishOption
	if len(opts) > 0 {
		opt = opts[0]
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return j.channel.Publish(
		opt.Exchange,
		j.queue,
		opt.Mandatory,
		opt.Immediate,
		amqp.Publishing{
			ContentType:  ContentTypeJson,
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (j *JsonQueueAdapter[T]) Close() {
	j.channel.Close()
	j.conn.Close()
}
