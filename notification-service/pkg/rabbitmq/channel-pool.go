package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
)

type ChannelPool struct {
	conn      *amqp.Connection
	pool      chan *amqp.Channel
	poolSize  int
	closeOnce sync.Once
}

func NewChannelPool(conn *amqp.Connection, poolSize int) (*ChannelPool, error) {
	pool := make(chan *amqp.Channel, poolSize)

	for i := 0; i < poolSize; i++ {
		ch, err := conn.Channel()
		if err != nil {
			return nil, fmt.Errorf("failed to create channel: %w", err)
		}
		pool <- ch
	}

	return &ChannelPool{
		conn:     conn,
		pool:     pool,
		poolSize: poolSize,
	}, nil
}

func (p *ChannelPool) Get() (*amqp.Channel, error) {
	select {
	case ch := <-p.pool:
		return ch, nil
	default:
		return p.conn.Channel()
	}
}

func (p *ChannelPool) Put(ch *amqp.Channel) {
	select {
	case p.pool <- ch:
	default:
		ch.Close()
	}
}

func (p *ChannelPool) Close() {
	p.closeOnce.Do(func() {
		close(p.pool)
		for ch := range p.pool {
			ch.Close()
		}
	})
}
