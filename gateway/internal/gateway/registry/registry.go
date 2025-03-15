package registry

import (
	"context"
	"fmt"
	"gateway/config"
	"gateway/internal/entity"
	"sync"
	"time"
)

type (
	serviceAddressesGetter interface {
		GetAllServicesByType(serviceType string) (services []string, err error)
	}
)

type Registry struct {
	serviceAddressesGetter serviceAddressesGetter
	services               sync.Map

	config config.Registry
}

func New(
	serviceAddressesGetter serviceAddressesGetter,
	config config.Registry,
) *Registry {
	return &Registry{
		serviceAddressesGetter: serviceAddressesGetter,
		services:               sync.Map{},
		config:                 config,
	}
}

func (r *Registry) RunActualizingRegistry(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Registry.Registry.RunActualizingRegistry: %w", err)
		}
	}()

	ticker := time.NewTicker(r.config.ActualizingInterval)
	for {
		select {
		case <-ticker.C:
			if err = r.Actualize(); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (r *Registry) GetAllWithType(serviceType string) (services []string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Registry.Registry.GetAllWithType: %w", err)
		}
	}()

	addresses, ok := r.services.Load(serviceType)
	if !ok {
		return nil, fmt.Errorf("service %s does not exist", serviceType)
	}

	convertedAddresses, ok := addresses.([]string)
	if !ok {
		return nil, fmt.Errorf("service %s is not of type []string", serviceType)
	}

	return convertedAddresses, nil
}

func (r *Registry) Actualize() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Registry.Registry.Actualize: %w", err)
		}
	}()

	newServicesMap := map[string][]string{}
	for _, serviceType := range entity.ServiceTypes {
		addresses, err := r.serviceAddressesGetter.GetAllServicesByType(serviceType)
		if err != nil {
			return err
		}

		newServicesMap[serviceType] = addresses
	}

	r.services.Range(func(key, value interface{}) bool {
		convertedKey, ok := key.(string)
		if !ok {
			return ok
		}

		r.services.Store(key, newServicesMap[convertedKey])

		return true
	})

	return nil
}
