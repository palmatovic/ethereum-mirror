package create

import (
	resource_perm_db "auth/pkg/database/resource_perm"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db           *gorm.DB
	resourcePerm *resource_perm_db.ResourcePerm
}

func NewService(db *gorm.DB, resourcePerm *resource_perm_db.ResourcePerm) *Service {
	return &Service{
		db:           db,
		resourcePerm: resourcePerm,
	}
}

func (s *Service) Create() (status int, group *resource_perm_db.ResourcePerm, err error) {
	if err = s.db.Create(s.resourcePerm).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fiber.StatusBadRequest, nil, err
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.resourcePerm, nil
}
