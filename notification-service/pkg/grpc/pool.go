package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strings"
	"sync"
	"time"
)

type (
	logger interface {
		Errorf(message string, args ...interface{})
	}

	registry interface {
		GetRandomServiceByType(serviceType string) (address string, err error)
	}
)

// Pool of the gRPC client connections
type Pool struct {
	pool *sync.Pool
}

// New creates new Pool
func New(logger logger, registry registry, serviceType string) *Pool {
	return &Pool{
		pool: &sync.Pool{
			New: func() interface{} {
				addr, err := registry.GetRandomServiceByType(serviceType)
				if err != nil {
					logger.Errorf("Pool.New: %s", err.Error())
					return err
				}
				newAddr := strings.Replace(strings.Split(addr, ":")[1], "/", "", -1)

				conn, err := grpc.NewClient(newAddr+":8081", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithIdleTimeout(time.Second))
				if err != nil {
					logger.Errorf("Pool.New: %s", err.Error())
					return err
				}

				return conn
			},
		},
	}
}

// Get returns grpc.ClientConn from the Pool
func (p *Pool) Get() *grpc.ClientConn {
	return p.pool.Get().(*grpc.ClientConn)
}

// Put put grpc.ClientConn back to the Pool
func (p *Pool) Put(client *grpc.ClientConn) {
	p.pool.Put(client)
}
