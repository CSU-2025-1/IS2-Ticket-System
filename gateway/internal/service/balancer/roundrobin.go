package balancer

import (
	"errors"
	"fmt"
	"math"
	"sync/atomic"
)

/*
	SUPER simplified implementation of round-robin algorithm without SMART features
*/

// RoundRobin is a classical load balancing method
type RoundRobin struct {
	registry               registry
	currentRoundRobinIndex atomic.Int64
}

// NewRoundRobin returns new example of RoundRobin
func NewRoundRobin(registry registry) *RoundRobin {
	return &RoundRobin{
		registry: registry,
	}
}

// GetAddress returns address of the service with RR algorithm
func (r *RoundRobin) GetAddress(serviceType string) (address string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Balancer.RoundRobin.GetAddress: %w", err)
		}
	}()
	defer r.currentRoundRobinIndex.Add(1)

	if r.currentRoundRobinIndex.Load() == math.MaxInt64 {
		r.currentRoundRobinIndex.Store(0)
	}

	addresses, err := r.registry.GetAllWithType(serviceType)
	if err != nil {
		return "", err
	}

	if len(addresses) == 0 {
		return "", errors.New("zero addresses received while balancing")
	}

	return addresses[(r.currentRoundRobinIndex.Load())%int64(len(addresses))], nil
}
