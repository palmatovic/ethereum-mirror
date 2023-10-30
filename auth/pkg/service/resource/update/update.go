package update

import (
	resource_db "auth/pkg/database/resource"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db       *gorm.DB
	resource *resource_db.Resource
}

func NewService(db *gorm.DB, resource *resource_db.Resource) *Service {
	return &Service{
		db:       db,
		resource: resource,
	}
}

func (s *Service) Update() (status int, resource *resource_db.Resource, err error) {
	if err = s.db.Where("ResourceId = ?", s.resource.ResourceId).Updates(s.resource).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.resource, nil
}
