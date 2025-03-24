package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

type (
	logger interface {
		Errorf(message string, args ...interface{})
	}

	registry interface {
		GetRandomServiceByType(serviceType string) (address string, err error)
	}
)

type Pool struct {
	pool *sync.Pool
}

func New(logger logger, registry registry, serviceType string) *Pool {
	return &Pool{
		pool: &sync.Pool{
			New: func() interface{} {
				addr, err := registry.GetRandomServiceByType(serviceType)
				if err != nil {
					logger.Errorf("Pool.New: %s", err.Error())
					return err
				}

				conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					logger.Errorf("Pool.New: %s", err.Error())
					return err
				}

				return conn
			},
		},
	}
}

func (p *Pool) Get() *grpc.ClientConn {
	return p.pool.Get().(*grpc.ClientConn)
}

func (p *Pool) Put(client *grpc.ClientConn) {
	p.pool.Put(client)
}
