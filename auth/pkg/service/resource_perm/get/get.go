package get

import (
	resource_perm_db "auth/pkg/database/resource_perm"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db             *gorm.DB
	resourcePermId int
}

func NewService(db *gorm.DB, resourcePermId int) *Service {
	return &Service{
		db:             db,
		resourcePermId: resourcePermId,
	}
}

func (s *Service) Get() (status int, group *resource_perm_db.ResourcePerm, err error) {
	group = new(resource_perm_db.ResourcePerm)
	if err = s.db.Where("ResourcePermRoleId = ?", s.resourcePermId).First(group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}
	return fiber.StatusOK, group, nil
}
