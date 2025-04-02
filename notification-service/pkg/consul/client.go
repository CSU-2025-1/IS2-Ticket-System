package consul

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"log/slog"
	"math/rand/v2"
	"strconv"
	"time"
)

// Client api wrapper for work with Consul
type Client struct {
	config     Config
	consul     *api.Client
	randomizer *rand.Rand
}

// New creates new Client
func New(config Config) *Client {
	return &Client{
		config: config,
		randomizer: rand.New(
			rand.NewPCG(
				uint64(time.Now().UnixMicro()),
				uint64(time.Now().UnixNano()),
			),
		),
	}
}

// Configure configures Consul client
func (c *Client) Configure() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Consul.Client.Configure: %w", err)
		}
	}()

	consul, err := api.NewClient(&api.Config{
		Address: c.config.ConsulAddress,
	})
	if err != nil {
		return err
	}

	c.consul = consul

	return nil
}

// GetAllServicesByType returns all service addresses with type = serviceType
func (c *Client) GetAllServicesByType(serviceType string) (addresses []string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Consul.Client.GetAllServicesByType: %w", err)
		}
	}()

	allServices, _, err := c.consul.
		Catalog().
		Services(
			&api.QueryOptions{
				Filter: fmt.Sprintf("\"%s\" in ServiceTags", serviceType),
			},
		)
	if err != nil {
		return nil, err
	}

	addresses = make([]string, 0, len(allServices))
	for serviceName := range allServices {
		entries, _, err := c.consul.Health().Service(serviceName, serviceType, true, nil)
		if err != nil {
			return nil, err
		}

		for entry := range entries {
			addresses = append(
				addresses,
				fmt.Sprintf("%s:%s", entries[entry].Service.Address, strconv.Itoa(entries[entry].Service.Port)),
			)
		}
	}

	return addresses, nil
}

// GetRandomServiceByType returns the random service address with type = serviceType
func (c *Client) GetRandomServiceByType(serviceType string) (address string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Consul.Client.GetRandomServiceByType: %w", err)
		}
	}()

	allServices, err := c.GetAllServicesByType(serviceType)
	if err != nil {
		return "", err
	}

	return allServices[c.randomizer.IntN(len(allServices))], nil
}

// GetFirstServiceByType returns the first found service address with type = serviceType
func (c *Client) GetFirstServiceByType(serviceType string) (address string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Consul.Client.GetFirstServiceByType: %w", err)
		}
	}()

	allServices, _, err := c.consul.
		Catalog().
		Services(
			&api.QueryOptions{
				Filter: fmt.Sprintf("\"%s\" in ServiceTags", serviceType),
			},
		)
	if err != nil {
		return "", err
	}

	for serviceName := range allServices {
		entries, _, err := c.consul.Health().Service(serviceName, serviceType, true, nil)
		if err != nil {
			return "", err
		}

		for entry := range entries {
			return fmt.Sprintf("%s:%s", entries[entry].Service.Address, strconv.Itoa(entries[entry].Service.Port)), nil
		}
	}

	return "", errors.New("service address not found")
}

// Register sends request for register service into Consul cluster
func (c *Client) Register(serviceType, address string, port uint16) (registeredServiceUUID string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Consul.Client.Register: %w", err)
		}
	}()

	registeredServiceUUID = uuid.NewString()

	registration := &api.AgentServiceRegistration{
		ID:      registeredServiceUUID,
		Name:    fmt.Sprintf("%s-%s", serviceType, address),
		Port:    int(port),
		Tags:    []string{serviceType},
		Address: fmt.Sprintf("http://%s", address),
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%v/check", address, port),
			Interval: c.config.HealthCheckInterval.String(),
			Timeout:  c.config.HealthCheckTimeout.String(),
		},
	}

	for {
		err := c.consul.Agent().ServiceRegister(registration)
		if err != nil {
			slog.Warn("consul connect failed, retry...")
			time.Sleep(2 * time.Second)
			continue
		}
		slog.Info("consul register success")
		return registeredServiceUUID, nil
	}
}

// Deregister sends request for deregister service into Consul cluster
func (c *Client) Deregister(serviceUUID string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Consul.Client.Deregister: %w", err)
		}
	}()

	return c.consul.Agent().ServiceDeregister(serviceUUID)
}
