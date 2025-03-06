package proxy

const (
	RoundRobin = "round-robin"
	Random     = "random"
)

type Proxy struct {
}

func NewProxy() *Proxy {
	return &Proxy{}
}
