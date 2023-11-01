package get_public_key

import (
	product_db "auth/pkg/database/product"
	"auth/pkg/service_util/aes"
	rsa_util "auth/pkg/service_util/rsa"
	"crypto/rsa"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db                  *gorm.DB
	name                string
	aes256EncryptionKey *aes.Key
}

func NewService(db *gorm.DB, aes256EncryptionKey *aes.Key, name string) *Service {
	return &Service{
		db:                  db,
		name:                name,
		aes256EncryptionKey: aes256EncryptionKey,
	}
}

func (s *Service) Get() (status int, publicKey *rsa.PublicKey, err error) {
	product := new(product_db.Product)
	if err = s.db.Where("Name = ?", s.name).First(product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}

	decryptedPublicKey, err := s.aes256EncryptionKey.Decrypt(product.RSAPublicKey)
	if err != nil {
		return 0, nil, err
	}

	publicKey, err = rsa_util.PublicKey(*decryptedPublicKey).ConvertToObj()
	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	return fiber.StatusOK, publicKey, nil
}
