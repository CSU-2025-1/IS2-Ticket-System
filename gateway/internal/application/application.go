package application

import (
	"context"
	"errors"
	"fmt"
	"gateway/config"
	"gateway/internal/repository"
	"gateway/internal/service/balancer"
	"gateway/internal/service/proxy"
	"gateway/internal/service/registry"
	"gateway/pkg/consul"
	"gateway/pkg/logger"
	"sync"
)

// Application is the main entrypoint for service project
type Application struct {
	config *config.Config
}

// New returns new Application with config
func New(config *config.Config) *Application {
	return &Application{
		config: config,
	}
}

// Run is a function that starting all application processes
func (a *Application) Run(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Application.Run: consul client configuring error: %w", err)
		}
	}()
	appLogger := logger.NewPrettyStdout(logger.Debug)
	auth := repository.NewAuthMock()
	consulClient := consul.New(
		consul.Config(a.config.Consul),
	)

	if err := consulClient.Configure(); err != nil {
		return err
	}

	serviceRegistry := registry.New(consulClient, appLogger, a.config.Registry)

	var gatewayProxy proxy.Proxy
	switch a.config.Proxy.BalancerAlgorithm {
	case balancer.RoundRobinAlgorithm:
		gatewayProxy = proxy.New(balancer.NewRoundRobin(serviceRegistry), auth, appLogger, a.config.Proxy)
	case balancer.RandomAlgorithm:
		gatewayProxy = proxy.New(balancer.NewRandom(serviceRegistry), auth, appLogger, a.config.Proxy)
	default:
		return errors.New(fmt.Sprintf("unknown load balancer algorithm: %s", a.config.Proxy.BalancerAlgorithm))
	}

	appLogger.Infof("Gateway starting...")
	wg := sync.WaitGroup{}
	go func() {
		defer wg.Done()

		wg.Add(1)
		if err := serviceRegistry.RunActualizingRegistry(ctx); err != nil {
			appLogger.Errorf("Application.Run: failed to run actualizing registry: %s", err.Error())
		}
	}()

	go func() {
		defer wg.Done()

		wg.Add(1)
		if err := gatewayProxy.Run(ctx); err != nil {
			appLogger.Errorf("Application.Run: failed to run proxy: %s", err.Error())
		}
	}()

	appLogger.Infof("Gateway is running now")
	wg.Wait()

	appLogger.Infof("Gateway stopped")
	return nil
}
