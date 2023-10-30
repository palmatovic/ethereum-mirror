package list

import (
	resource_db "auth/pkg/database/resource"
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

func (s *Service) List() (status int, resources *[]resource_db.Resource, err error) {
	resources = new([]resource_db.Resource)
	if err = s.db.Order("Name ASC").Offset(s.offset).Limit(s.limit).Find(resources).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, resources, nil
}
