package get_server_ssl

import (
	product_db "auth/pkg/database/product"
	"auth/pkg/service_util/aes"
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

func (s *Service) Get() (status int, serverSslCert []byte, serverSslKey []byte, err error) {
	product := new(product_db.Product)
	if err = s.db.Where("Name = ?", s.name).First(product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, nil, err
		}
	}

	decryptedServerCert, err := s.aes256EncryptionKey.Decrypt(product.ServerCert)
	if err != nil {
		return fiber.StatusInternalServerError, nil, nil, err
	}

	decryptedServerKey, err := s.aes256EncryptionKey.Decrypt(product.ServerKey)
	if err != nil {
		return fiber.StatusInternalServerError, nil, nil, err
	}

	return fiber.StatusOK, []byte(*decryptedServerCert), []byte(*decryptedServerKey), nil
}
