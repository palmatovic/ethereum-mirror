package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
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

type PublicKey []byte

func (p PublicKey) ConvertToObj() (*rsa.PublicKey, error) {
	block, _ := pem.Decode(p)
	if block == nil {
		return nil, errors.New("error during decoding public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPublicKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("rsa.PublicKey is not a *rsa.PublicKey")
	}
	return rsaPublicKey, nil
}
