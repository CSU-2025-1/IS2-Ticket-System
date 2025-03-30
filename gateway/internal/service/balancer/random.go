package balancer

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// Random is a classical load balancing method with random address choosing
type Random struct {
	registry   registry
	randomizer *rand.Rand
}

// NewRandom returns new example of Random
func NewRandom(registry registry) *Random {
	return &Random{
		registry: registry,
		randomizer: rand.New(
			rand.NewPCG(
				uint64(time.Now().UnixMicro()),
				uint64(time.Now().UnixNano()),
			),
		),
	}
}

// GetAddress returns address of the random service
func (r *Random) GetAddress(serviceType string) (address string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Balancer.Random.GetAddress: %w", err)
		}
	}()

	addresses, err := r.registry.GetAllWithType(serviceType)
	if err != nil {
		return "", err
	}

	return addresses[r.randomizer.Uint64()%uint64(len(address))], nil
}
