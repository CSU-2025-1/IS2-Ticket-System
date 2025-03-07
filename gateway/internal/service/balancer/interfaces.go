package balancer

import (
	"context"
	"gateway/internal/entity"
)

type (
	registry interface {
		GetAllWithType(ctx context.Context, serviceType string) ([]entity.Service, error)
	}
)
