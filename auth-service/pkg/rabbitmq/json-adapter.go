package rabbitmq

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"sync"
)

type JsonQueueAdapter[T any] struct {
	pool      *ChannelPool
	queue     string
	wg        sync.WaitGroup
	closeOnce sync.Once
}

func NewJsonQueueAdapter[T any](pool *ChannelPool, queue string) *JsonQueueAdapter[T] {
	return &JsonQueueAdapter[T]{
		pool:  pool,
		queue: queue,
	}
}

func (j *JsonQueueAdapter[T]) Consume(opts ...ConsumeOption) (<-chan T, error) {
	ch, err := j.pool.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}

	opt := DefaultConsumeOption
	if len(opts) > 0 {
		opt = opts[0]
	}

	msgs, err := ch.Consume(
		j.queue,
		opt.Consumer,
		opt.AutoAck,
		opt.Exclusive,
		opt.NoLocal,
		opt.NoWait,
		opt.Args,
	)
	if err != nil {
		j.pool.Put(ch)
		return nil, fmt.Errorf("failed to consume: %w", err)
	}

	output := make(chan T)
	j.wg.Add(1)

	go func() {
		defer j.wg.Done()
		defer close(output)
		defer j.pool.Put(ch)

		for d := range msgs {
			var msg T
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				slog.Error("failed to unmarshal json message", "error", err)
				continue
			}
			output <- msg
			if !opt.AutoAck {
				d.Ack(false)
			}
		}
	}()

	return output, nil
}

func (j *JsonQueueAdapter[T]) Publish(msg T, opts ...PublishOption) error {
	ch, err := j.pool.Get()
	if err != nil {
		return fmt.Errorf("failed to get channel: %w", err)
	}
	defer j.pool.Put(ch)

	opt := DefaultPublishOption
	if len(opts) > 0 {
		opt = opts[0]
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return ch.Publish(
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
	j.closeOnce.Do(func() {
		j.wg.Wait()
		j.pool.Close()
	})
}
