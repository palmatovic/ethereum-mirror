package update

import (
	product_db "auth/pkg/database/product"
	model_update_product "auth/pkg/model/api/product/update"
	"auth/pkg/service/product/get"
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
	product             *model_update_product.Product
	aes256EncryptionKey *aes.Key
}

func NewService(db *gorm.DB, product *model_update_product.Product, aes256EncryptionKey *aes.Key) *Service {
	return &Service{
		db:                  db,
		product:             product,
		aes256EncryptionKey: aes256EncryptionKey,
	}
}

func (s *Service) Update() (status int, product *product_db.Product, err error) {
	findStatus, dbProduct, err := get.NewService(s.db, s.product.ProductId).Get()
	if err != nil {
		if findStatus == fiber.StatusNotFound {
			return fiber.StatusNotFound, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}

	var count int
	if s.product.RenewSsl != nil {
		count++
	}
	if s.product.RenewRSA256 != nil {
		count++
	}

	if count != 1 {
		return fiber.StatusBadRequest, nil, errors.New("only one of renew_ssl or renew_rsa_256 can be set")
	}

	switch {
	case s.product.RenewSsl != nil:
		certificates, err := ssl.NewService(
			s.product.RenewSsl.Company,
			s.product.RenewSsl.Country,
			s.product.RenewSsl.Province,
			s.product.RenewSsl.Locality,
			s.product.RenewSsl.CompanyUnit,
			s.product.RenewSsl.CommonName,
			s.product.RenewSsl.AltDNS,
		).NewCertificates()
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

		dbProduct.ServerCert = *encryptedServerCrt
		dbProduct.ServerKey = *encryptedServerKey
		dbProduct.CaCert = *encryptedCACrt
		dbProduct.CaKey = *encryptedCAKey
		dbProduct.SSLExpired = false
		dbProduct.SSLExpirationDate = time.Now().Add(time.Hour * 24 * 365)
		break
	case s.product.RenewRSA256 != nil:
		if s.product.RenewRSA256.RenewKeyPair != nil && *s.product.RenewRSA256.RenewKeyPair == true {
			keys, err := rsa.NewRsa().GenerateRSAKeys()
			if err != nil {
				return fiber.StatusInternalServerError, nil, err
			}

			encryptedPrivateKey, err := s.aes256EncryptionKey.Encrypt(string(keys.Private))
			if err != nil {
				return 0, nil, err
			}
			encryptedPublickey, err := s.aes256EncryptionKey.Encrypt(string(keys.Public))
			if err != nil {
				return 0, nil, err
			}

			dbProduct.RSAPrivateKey = *encryptedPrivateKey
			dbProduct.RSAPublicKey = *encryptedPublickey
			dbProduct.RSAExpirationDate = time.Now().Add(365 * 24 * time.Hour)
			dbProduct.RSAExpired = false
		}
		if s.product.RenewRSA256.RenewConfig != nil {
			dbProduct.AccessTokenExpiresInMinutes = s.product.RenewRSA256.RenewConfig.AccessToken.ExpiresInMinutes
			dbProduct.RefreshTokenExpiresInMinutes = s.product.RenewRSA256.RenewConfig.RefreshToken.ExpiresInMinutes
		}
		break
	}

	if err = s.db.Where("ProductId = ?", dbProduct.ProductId).Updates(dbProduct).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, dbProduct, nil
}
