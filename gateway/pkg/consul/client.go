package consul

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"net"
	"strconv"
)

// Client api wrapper for work with Consul
type Client struct {
	config *api.Config
	consul *api.Client
}

// New creates new Client
func New(config Config) Client {
	return Client{
		config: &api.Config{
			Address: config.Address,
		},
	}
}

// Configure configures Consul client
func (c *Client) Configure() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Consul.Client.Configure: %w", err)
		}
	}()

	if c.config == nil {
		return errors.New("client not initialized")
	}

	consul, err := api.NewClient(c.config)
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
				net.JoinHostPort(
					entries[entry].Service.Address,
					strconv.Itoa(entries[entry].Service.Port),
				),
			)
		}
	}

	return addresses, nil
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
			return net.JoinHostPort(
				entries[entry].Service.Address,
				strconv.Itoa(entries[entry].Service.Port),
			), nil
		}
	}

	return "", errors.New("service address not found")
}
