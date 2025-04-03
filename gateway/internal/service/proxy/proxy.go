package proxy

import (
	"context"
	"errors"
	"fmt"
	"gateway/config"
	"net"
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
		Auth(ctx context.Context, token string) (map[string]string, error)
	}

	logger interface {
		Infof(format string, args ...interface{})
		Warnf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
	}
)

// Proxy is the HTTP v1 proxy for request to services of Ticket system
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
) *Proxy {
	return &Proxy{
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
	mux.HandleFunc("/", p.handle)
	server := &http.Server{
		BaseContext: func(net.Listener) context.Context { return ctx },
		Addr:        fmt.Sprintf(":%d", p.config.LaunchedPort),
		Handler:     mux,
	}

	p.logger.Infof("Proxy.Run: server is running")

	if err = server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			p.logger.Infof("Proxy.Run: http server closed")
			return nil
		}

		return err
	}

	if err := server.Close(); err != nil {
		p.logger.Errorf("Proxy.Run: server closed with error: %s", err.Error())
	}

	p.logger.Infof("Proxy.Run: server stopped")

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

	p.logger.Infof("Proxy.handle: handled request with path: %s and proxied to: %s", requestPath, newUrl.String())
	proxy := httputil.NewSingleHostReverseProxy(newUrl)
	request.URL.Scheme = newUrl.Scheme
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

	return p.authenticator.Auth(request.Context(), parts[1])
}
