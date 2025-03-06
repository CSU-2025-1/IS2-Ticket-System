package cache

import (
	"gateway/internal/entity"
	"sync"
)

type InMemoryEntity struct {
	cache map[string][]entity.Service
	mutex sync.RWMutex
}

func NewInMemoryEntity() *InMemoryEntity {
	return &InMemoryEntity{
		cache: make(map[string][]entity.Service),
		mutex: sync.RWMutex{},
	}
}

func (i *InMemoryEntity) GetWithType(serviceType string) []entity.Service {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	return i.cache[serviceType]
}

func (i *InMemoryEntity) Add(service entity.Service) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	services := i.cache[service.Type]
	services = append(services, service)
	i.cache[service.Type] = services
}
