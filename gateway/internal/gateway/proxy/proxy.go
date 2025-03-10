package proxy

import "gateway/internal/config"

type (
	balancer interface {
		GetAddress(serviceType string) (string, error)
	}

	authenticator interface {
		Auth() (map[string]string, error)
	}
)

type Proxy struct {
	balancer      balancer
	authenticator authenticator
	config        config.Proxy
}

func New(
	balancer balancer,
	authenticator authenticator,
	config config.Proxy,
) Proxy {
	return Proxy{
		balancer:      balancer,
		authenticator: authenticator,
		config:        config,
	}
}

func (p *Proxy) ListenAndServe() {

}
