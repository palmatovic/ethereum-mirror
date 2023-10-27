package get

import (
	company_db "auth/pkg/database/company"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	companyId int
}

func NewService(db *gorm.DB, companyId int) *Service {
	return &Service{
		db:        db,
		companyId: companyId,
	}
}

func (s *Service) Get() (status int, company *company_db.Company, err error) {
	company = new(company_db.Company)
	if err = s.db.Where("CompanyId = ?", s.companyId).First(company).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}
	return fiber.StatusOK, company, nil
}
