package hasher

import "crypto/sha512"

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
		res = append(res, s.salt...)
		res = hasher.Sum(res)
	}
	return string(res)
}
