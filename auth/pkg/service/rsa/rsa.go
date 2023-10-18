package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type Rsa struct {
	baseName string
}

type KeyPair struct {
	PrivateKey string
	PublicKey  string
}

func NewRsa(baseName string) *Rsa {
	return &Rsa{
		baseName: baseName,
	}
}

func (r *Rsa) GenerateRSAKeys() (*KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	privateKeyPEMString := string(pem.EncodeToMemory(privateKeyPEM))

	publicKeyPEM, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	publicKeyPEMString := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyPEM}))

	return &KeyPair{
		PrivateKey: privateKeyPEMString,
		PublicKey:  publicKeyPEMString,
	}, nil
}
