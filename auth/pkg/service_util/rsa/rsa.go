package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type Rsa struct{}
type Key []byte

type Pair struct {
	Private Key
	Public  Key
}

func NewRsa() *Rsa {
	return &Rsa{}
}

func (r *Rsa) GenerateRSAKeys() (*Pair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	publicKeyPEM, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	return &Pair{
		Private: pem.EncodeToMemory(privateKeyPEM),
		Public:  pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyPEM}),
	}, nil
}
