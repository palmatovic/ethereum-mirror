package update

import (
	company_db "auth/pkg/database/company"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db      *gorm.DB
	company *company_db.Company
}

func NewService(db *gorm.DB, company *company_db.Company) *Service {
	return &Service{
		db:      db,
		company: company,
	}
}

func (s *Service) Update() (status int, company *company_db.Company, err error) {
	if err = s.db.Where("CompanyId = ?", s.company.CompanyId).Updates(s.company).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.company, nil
}
