package delete

import (
	resource_db "auth/pkg/database/resource"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	resourceId int64
}

func NewService(db *gorm.DB, resourceId int64) *Service {
	return &Service{
		db:         db,
		resourceId: resourceId,
	}
}

func (s *Service) Delete() (status int, resource *resource_db.Resource, err error) {
	resource = new(resource_db.Resource)
	if err = s.db.Where("ResourceId = ?", s.resourceId).Delete(resource).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) || errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, resource, nil
}
