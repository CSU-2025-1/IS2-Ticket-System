package entity

import (
	"fmt"
	"net"
)

type Service struct {
	Type string
	IP   net.IP
	Port uint32
}

func (s *Service) GetFullAddress() string {
	return fmt.Sprintf("%s:%d", s.IP.String(), s.Port)
}
