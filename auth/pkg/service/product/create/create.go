package create

import (
	product_db "auth/pkg/database/product"
	model_create_product "auth/pkg/model/api/product/create"
	"auth/pkg/service_util/rsa"
	"auth/pkg/service_util/ssl"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

type Service struct {
	db      *gorm.DB
	product *model_create_product.Product
}

func NewService(db *gorm.DB, product *model_create_product.Product) *Service {
	return &Service{
		db:      db,
		product: product,
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

	newProduct = &product_db.Product{
		Name:                         s.product.Name,
		Description:                  s.product.Description,
		ServerCert:                   certificates.Server.Cert,
		ServerKey:                    certificates.Server.Key,
		CaCert:                       certificates.CA.Cert,
		CaKey:                        certificates.CA.Key,
		SSLExpired:                   false,
		RSAPrivateKey:                keys.Private,
		RSAPublicKey:                 keys.Public,
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
