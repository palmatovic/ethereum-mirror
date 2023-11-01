package sha

import (
	"crypto/sha256"
	"encoding/hex"
)

type Service struct {
	input string
}

func NewService(input string) *Service {
	return &Service{input: input}
}

func (s *Service) Sha256() string {
	h := sha256.New()
	h.Write([]byte(s.input))
	return hex.EncodeToString(h.Sum(nil))
}
