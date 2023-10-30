package list

import (
	company_db "auth/pkg/database/company"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	limit  int64
	offset int64
}

func NewService(db *gorm.DB, pageSize int64, pageNumber int64) *Service {
	return &Service{
		db:     db,
		limit:  pageSize,
		offset: (pageNumber - 1) * pageSize,
	}
}

func (s *Service) List() (status int, companies *[]company_db.Company, err error) {
	companies = new([]company_db.Company)
	if err = s.db.Order("Name ASC").Offset(int(s.offset)).Limit(int(s.limit)).Find(companies).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, companies, nil
}
