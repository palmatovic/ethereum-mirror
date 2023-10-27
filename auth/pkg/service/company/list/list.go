package list

import (
	company_db "auth/pkg/database/company"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	limit  int
	offset int
}

func NewService(db *gorm.DB, pageSize int, pageNumber int) *Service {
	return &Service{
		db:     db,
		limit:  pageSize,
		offset: (pageNumber - 1) * pageSize,
	}
}

func (s *Service) List() (status int, companys *[]company_db.Company, err error) {
	companys = new([]company_db.Company)
	if err = s.db.Order("Name ASC").Offset(s.offset).Limit(s.limit).Find(companys).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, companys, nil
}
