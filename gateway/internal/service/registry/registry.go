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

	logger interface {
		Debugf(format string, args ...interface{})
		Infof(format string, args ...interface{})
	}
)

type Registry struct {
	serviceAddressesGetter serviceAddressesGetter
	logger                 logger
	services               sync.Map
	config                 config.Registry
}

func New(
	serviceAddressesGetter serviceAddressesGetter,
	logger logger,
	config config.Registry,
) *Registry {
	return &Registry{
		serviceAddressesGetter: serviceAddressesGetter,
		logger:                 logger,
		services:               sync.Map{},
		config:                 config,
	}
}

// RunActualizingRegistry actualize internal service list from external registry
func (r *Registry) RunActualizingRegistry(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Registry.RunActualizingRegistry: %w", err)
		}
	}()

	r.logger.Infof("Registry.RunActualizingRegistry: registry actualizing is running")

	ticker := time.NewTicker(r.config.ActualizingInterval)
	for {
		select {
		case <-ticker.C:
			r.logger.Infof("Registry.RunActualizingRegistry: actualizing registry procceed...")
			if err = r.Actualize(); err != nil {
				return err
			}
		case <-ctx.Done():
			r.logger.Infof("Registry.RunActualizingRegistry: stopped")
			return nil
		}
	}
}

func (r *Registry) GetAllWithType(serviceType string) (services []string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Registry.GetAllWithType: %w", err)
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

// Actualize is a function for filling services from registry into internal gateway service list
func (r *Registry) Actualize() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Registry.Actualize: %w", err)
		}
	}()

	for _, serviceType := range entity.ServiceTypes {
		addresses, err := r.serviceAddressesGetter.GetAllServicesByType(serviceType)
		if err != nil {
			return err
		}

		r.logger.Debugf("Registry.Actualize: get %d service addresses with type: %s and these addresses: %v",
			len(addresses),
			serviceType,
			addresses,
		)

		r.services.Store(serviceType, addresses)
	}

	return nil
}
