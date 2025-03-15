package proxy

import (
	"context"
	"errors"
	"fmt"
	"gateway/config"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type (
	balancer interface {
		GetAddress(serviceType string) (string, error)
	}

	authenticator interface {
		Auth(token string) (map[string]string, error)
	}

	logger interface {
		Infof(format string, args ...interface{})
		Warnf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
	}
)

type Proxy struct {
	balancer      balancer
	authenticator authenticator
	logger        logger
	config        config.Proxy
}

func New(
	balancer balancer,
	authenticator authenticator,
	logger logger,
	config config.Proxy,
) Proxy {
	return Proxy{
		balancer:      balancer,
		authenticator: authenticator,
		logger:        logger,
		config:        config,
	}
}

func (p *Proxy) Run(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Proxy.Run: %w", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	if err = http.ListenAndServe(fmt.Sprintf(":%d", p.config.LaunchedPort), mux); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			p.logger.Infof("Proxy.Run: http server closed")
			return nil
		}

		return err
	}

	return nil
}

func (p *Proxy) handle(responseWriter http.ResponseWriter, request *http.Request) {
	requestPath := request.URL.Path

	newHeadersData := map[string]string{}
	if p.IsForAuthorized(requestPath) {
		authData, err := p.Authenticate(request)
		if err != nil {
			p.logger.Infof("Proxy.handle: failed to authenticate: %s", err.Error())
			http.Error(responseWriter, err.Error(), http.StatusForbidden)
			return
		}

		newHeadersData = authData
	}

	for header, value := range newHeadersData {
		request.Header.Set(header, value)
	}

	address, err := p.MatchRouteWithService(requestPath)
	if err != nil {
		p.logger.Infof("Proxy.handle: failed to get address from balancer: %s", err.Error())
		http.Error(responseWriter, "Critical error with load balancing", http.StatusBadGateway)
		return
	}

	newUrl, err := url.Parse(address)
	if err != nil {
		p.logger.Infof("Proxy.handle: failed to parse url: %s", err.Error())
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(newUrl)
	request.URL.Host = request.Host
	request.Host = newUrl.Host
	proxy.ServeHTTP(responseWriter, request)
}

func (p *Proxy) MatchRouteWithService(requestPath string) (address string, err error) {
	for path, serviceType := range RoutePolicy {
		if strings.HasPrefix(requestPath, path) {
			return p.balancer.GetAddress(serviceType)
		}
	}

	return "", errors.New("no route policy found and matched")
}

func (p *Proxy) IsForAuthorized(requestPath string) bool {
	for path, isForAuthorized := range AuthPolicy {
		if strings.HasPrefix(requestPath, path) {
			return isForAuthorized
		}
	}

	return false
}

func (p *Proxy) Authenticate(request *http.Request) (authData map[string]string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Proxy.Authenticate: %w", err)
		}
	}()

	header := request.Header.Get("Authorization")
	if header == "" {
		return nil, errors.New("no authorization header")
	}

	parts := strings.Split(header, " ")
	if len(parts) < 2 {
		return nil, errors.New("invalid authorization header")
	}

	if parts[0] != "Bearer" {
		return nil, errors.New("not a Bearer token")
	}

	return p.authenticator.Auth(parts[1])
}
