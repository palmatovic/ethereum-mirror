package aes

import (
	"crypto/rand"
	"io"
)

type Service struct{}
type Key []byte

func NewService() *Service { return &Service{} }

func (s *Service) NewAES256Key() (key *Key, err error) {
	key = new(Key)
	*key = make([]byte, 32)
	if _, err = io.ReadFull(rand.Reader, *key); err != nil {
		return nil, err
	}
	return key, nil
}
