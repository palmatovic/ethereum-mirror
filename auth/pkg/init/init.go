package init

import (
	product_db "auth/pkg/database/product"
	"auth/pkg/model/api/product/create"
	product_create_service "auth/pkg/service/product/create"
	product_get_by_name_service "auth/pkg/service/product/get_by_name"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db}
}

func (s *Service) Init() {
	var err error
	if _, _, err = product_get_by_name_service.NewService(s.db, "auth").Get(); err == nil {
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Panic("terminated with failure")
	}

	var product *product_db.Product
	if _, product, err = product_create_service.NewService(s.db, &create.Product{
		Name:        "auth",
		Description: "auth product",
		SslSetup: create.SslSetup{
			Company:     "auth-company",
			Province:    "Rome",
			Country:     "IT",
			CompanyUnit: "ENG",
			CommonName:  "*.auth-company.com",
			Locality:    "Rome",
			AltDNS:      "*.auth-company.com",
		},
		JwtConfig: create.JwtConfig{
			AccessToken:  create.Token{ExpiresInMinutes: int64(time.Minute * 4 * 60)},
			RefreshToken: create.Token{ExpiresInMinutes: int64(time.Minute * 8 * 60)},
		},
	}).Create(); err != nil {
		logrus.WithError(err).Panic("terminated with failure")
	}

	// creare gruppi, ruoli, risorse ecc

}
