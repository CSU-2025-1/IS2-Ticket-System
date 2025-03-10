package balancer

type (
	registry interface {
		GetAllWithType(serviceType string) (addresses []string, err error)
	}
)
