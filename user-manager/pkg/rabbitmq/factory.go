package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
)

var (
	globalPool     *ChannelPool
	globalPoolOnce sync.Once
)

func InitGlobalPool(cfg Config) (*ChannelPool, error) {
	var initErr error
	globalPoolOnce.Do(func() {
		conn, err := amqp.Dial(cfg.ToDSN())
		if err != nil {
			initErr = fmt.Errorf("failed to connect: %w", err)
			return
		}

		pool, err := NewChannelPool(conn, cfg.PoolSize)
		if err != nil {
			conn.Close()
			initErr = fmt.Errorf("failed to create pool: %w", err)
			return
		}

		globalPool = pool
	})
	return globalPool, initErr
}

func CreateJsonWriter[T any](cfg Config, option PublishOption) (*Writer[T], error) {
	pool, err := InitGlobalPool(cfg)
	if err != nil {
		return nil, err
	}
	return &Writer[T]{
		adapter: NewJsonQueueAdapter[T](pool, cfg.Queue),
		option:  option,
	}, nil
}

func CreateJsonReader[T any](cfg Config, handler func(T), opts ConsumeOption) (*Reader[T], error) {
	pool, err := InitGlobalPool(cfg)
	if err != nil {
		return nil, err
	}
	return &Reader[T]{
		adapter: NewJsonQueueAdapter[T](pool, cfg.Queue),
		handler: handler,
		option:  opts,
	}, nil
}
