package create

import (
	product_db "auth/pkg/database/product"
	model_create_product "auth/pkg/model/api/product/create"
	"auth/pkg/service_util/aes"
	"auth/pkg/service_util/rsa"
	"auth/pkg/service_util/ssl"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

type Service struct {
	db                  *gorm.DB
	product             *model_create_product.Product
	aes256EncryptionKey *aes.Key
}

func NewService(db *gorm.DB, product *model_create_product.Product, aes256EncryptionKey *aes.Key) *Service {
	return &Service{
		db:                  db,
		product:             product,
		aes256EncryptionKey: aes256EncryptionKey,
	}
}

func (s *Service) Create() (status int, newProduct *product_db.Product, err error) {
	certificates, err := ssl.NewService(
		s.product.SslSetup.Company,
		s.product.SslSetup.Country,
		s.product.SslSetup.Province,
		s.product.SslSetup.Locality,
		s.product.SslSetup.CompanyUnit,
		s.product.SslSetup.CommonName,
		s.product.SslSetup.AltDNS,
	).NewCertificates()
	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	keys, err := rsa.NewRsa().GenerateRSAKeys()
	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	encryptedCAKey, err := s.aes256EncryptionKey.Encrypt(string(certificates.CA.Key))
	if err != nil {
		return 0, nil, err
	}
	encryptedCACrt, err := s.aes256EncryptionKey.Encrypt(string(certificates.CA.Cert))
	if err != nil {
		return 0, nil, err
	}
	encryptedServerKey, err := s.aes256EncryptionKey.Encrypt(string(certificates.Server.Key))
	if err != nil {
		return 0, nil, err
	}
	encryptedServerCrt, err := s.aes256EncryptionKey.Encrypt(string(certificates.Server.Cert))
	if err != nil {
		return 0, nil, err
	}

	encryptedRsaPrivateKey, err := s.aes256EncryptionKey.Encrypt(string(keys.Private))
	if err != nil {
		return 0, nil, err
	}
	encryptedRsaPublicKey, err := s.aes256EncryptionKey.Encrypt(string(keys.Public))
	if err != nil {
		return 0, nil, err
	}

	newProduct = &product_db.Product{
		Name:                         s.product.Name,
		Description:                  s.product.Description,
		ServerCert:                   *encryptedServerCrt,
		ServerKey:                    *encryptedServerKey,
		CaCert:                       *encryptedCACrt,
		CaKey:                        *encryptedCAKey,
		SSLExpired:                   false,
		RSAPrivateKey:                *encryptedRsaPrivateKey,
		RSAPublicKey:                 *encryptedRsaPublicKey,
		AccessTokenExpiresInMinutes:  s.product.JwtConfig.AccessToken.ExpiresInMinutes,
		RefreshTokenExpiresInMinutes: s.product.JwtConfig.RefreshToken.ExpiresInMinutes,
		RSAExpirationDate:            time.Now().Add(365 * 24 * time.Hour),
		RSAExpired:                   false,
	}

	if err = s.db.Create(newProduct).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, newProduct, nil
}
