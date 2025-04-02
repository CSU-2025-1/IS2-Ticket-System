package hasher

import (
	"crypto/sha512"
	"fmt"
)

type Service struct {
	cost int
	salt []byte
}

func New(cost int, salt string) *Service {
	return &Service{
		cost: cost,
		salt: []byte(salt),
	}
}

func (s *Service) Hash(in string) string {
	hasher := sha512.New()
	bytes := []byte(in)
	res := bytes
	for range s.cost {
		hasher.Write(res)
		res = hasher.Sum(nil)
	}
	return fmt.Sprintf("%x", res)
}
