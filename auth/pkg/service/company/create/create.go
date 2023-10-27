package create

import (
	company_db "auth/pkg/database/company"
	"errors"
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

func (s *Service) Create() (status int, company *company_db.Company, err error) {
	if err = s.db.Create(s.company).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.company, nil
}
