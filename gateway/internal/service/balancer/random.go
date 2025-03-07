package balancer

import (
	"context"
	"gateway/internal/entity"
	"math/rand/v2"
	"time"
)

const (
	randomGetOp = "RandomBalancer.Get"
)

type RandomBalancer struct {
	registry registry
	rand     *rand.Rand
}

func NewRandomBalancer(registry registry) RandomBalancer {
	return RandomBalancer{
		registry: registry,
		rand: rand.New(
			rand.NewPCG(
				uint64(time.Now().UnixMicro()),
				uint64(time.Now().UnixMicro()),
			),
		),
	}
}

func (r *RandomBalancer) Get(serviceType string) (entity.Service, error) {
	_, _ = r.registry.GetAllWithType(context.Background(), serviceType)

	panic("implement me")
}
