package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
)

type ConnectionPool struct {
	config Config
	pool   chan *amqp.Connection
	mu     sync.Mutex
	closed bool
}

func NewConnectionPool(cfg Config) (*ConnectionPool, error) {
	pool := make(chan *amqp.Connection, cfg.PoolSize)

	for i := 0; i < cfg.PoolSize; i++ {
		conn, err := amqp.Dial(cfg.ToDSN())
		if err != nil {
			return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
		}
		pool <- conn
	}

	return &ConnectionPool{
		config: cfg,
		pool:   pool,
	}, nil
}

func (p *ConnectionPool) Get() (*amqp.Connection, error) {
	select {
	case conn := <-p.pool:
		return conn, nil
	default:
		return amqp.Dial(p.config.ToDSN())
	}
}

func (p *ConnectionPool) Put(conn *amqp.Connection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		conn.Close()
		return
	}

	select {
	case p.pool <- conn:
	default:
		conn.Close()
	}
}

func (p *ConnectionPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	p.closed = true
	close(p.pool)

	for conn := range p.pool {
		conn.Close()
	}
}
