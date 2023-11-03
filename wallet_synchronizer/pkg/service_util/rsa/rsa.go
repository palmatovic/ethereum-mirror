package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

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

type PublicKeyFilepath string

func (p PublicKeyFilepath) ConvertToObj() (*rsa.PublicKey, error) {
	file, err := os.ReadFile(string(p))
	if err != nil {
		return nil, err
	}
	return PublicKey(file).ConvertToObj()
}
