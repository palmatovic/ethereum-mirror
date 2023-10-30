package create

import (
	resource_db "auth/pkg/database/resource"
	"errors"
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

func (s *Service) Create() (status int, resource *resource_db.Resource, err error) {
	if err = s.db.Create(s.resource).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fiber.StatusBadRequest, nil, err
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.resource, nil
}
