package update

import (
	product_db "auth/pkg/database/product"
	model_update_product "auth/pkg/model/api/product/update"
	"auth/pkg/service/product/get"
	"auth/pkg/service_util/rsa"
	"auth/pkg/service_util/ssl"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

type Service struct {
	db      *gorm.DB
	product *model_update_product.Product
}

func NewService(db *gorm.DB, product *model_update_product.Product) *Service {
	return &Service{
		db:      db,
		product: product,
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

	if s.product.SslSetup != nil {
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
		
		dbProduct.ServerCert = certificates.Server.Cert
		dbProduct.ServerKey = certificates.Server.Key
		dbProduct.CaCert = certificates.CA.Cert
		dbProduct.CaKey = certificates.CA.Key
		dbProduct.SSLExpired = false
	}

	if s.product.JwtConfig != nil {
		if s.product.JwtConfig.RenewRSA256 {
			keys, err := rsa.NewRsa().GenerateRSAKeys()
			if err != nil {
				return fiber.StatusInternalServerError, nil, err
			}
			dbProduct.RSAPrivateKey = keys.Private
			dbProduct.RSAPublicKey = keys.Public
			dbProduct.RSAExpirationDate = time.Now().Add(365 * 24 * time.Hour)
			dbProduct.RSAExpired = false
		}
		if s.product.JwtConfig.RenewTokenConfig != nil {
			dbProduct.AccessTokenExpiresInMinutes = s.product.JwtConfig.RenewTokenConfig.AccessToken.ExpiresInMinutes
			dbProduct.RefreshTokenExpiresInMinutes = s.product.JwtConfig.RenewTokenConfig.RefreshToken.ExpiresInMinutes
		}
	}

	if err = s.db.Where("ProductId = ?", dbProduct.ProductId).Updates(dbProduct).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, dbProduct, nil
}
