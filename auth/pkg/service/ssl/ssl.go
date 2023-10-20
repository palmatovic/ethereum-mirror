package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

type Service struct {
	organization     string
	country          string
	province         string
	locality         string
	organizationUnit string
	commonName       string
	altDNS           string
}

func NewService(
	organization string,
	country string,
	province string,
	locality string,
	organizationUnit string,
	commonName string,
	altDNS string,
) *Service {
	return &Service{
		organization:     organization,
		country:          country,
		province:         province,
		locality:         locality,
		organizationUnit: organizationUnit,
		commonName:       commonName,
		altDNS:           altDNS,
	}
}

type Certificates struct {
	Server struct {
		Cert []byte
		Key  []byte
	}
	CA struct {
		Cert []byte
		Key  []byte
	}
}

func (s *Service) NewCertificates() (*Certificates, error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Organization:       []string{s.organization},
			Country:            []string{s.country},
			Province:           []string{s.province},
			Locality:           []string{s.locality},
			OrganizationalUnit: []string{s.organizationUnit},
			CommonName:         s.commonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create our private and public key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, err
	}

	// pem encode
	caPEM := new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return nil, err
	}

	caPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	if err != nil {
		return nil, err
	}

	// set up our server certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Organization:       []string{s.organization},
			Country:            []string{s.country},
			Province:           []string{s.province},
			Locality:           []string{s.locality},
			OrganizationalUnit: []string{s.organizationUnit},
			CommonName:         s.commonName,
		},
		DNSNames:    []string{s.altDNS},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(10, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, err
	}

	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return nil, err
	}

	certPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	if err != nil {
		return nil, err
	}

	return &Certificates{
		Server: struct {
			Cert []byte
			Key  []byte
		}{certPEM.Bytes(), certPrivKeyPEM.Bytes()},
		CA: struct {
			Cert []byte
			Key  []byte
		}{caPEM.Bytes(), caPrivKeyPEM.Bytes()},
	}, nil
}
