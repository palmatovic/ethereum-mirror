package update

import (
	resource_perm_db "auth/pkg/database/resource_perm"
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

func (s *Service) Update() (status int, group *resource_perm_db.ResourcePerm, err error) {
	if err = s.db.Where("ResourcePermRoleId = ?", s.resourcePerm.ResourcePermId).Updates(s.resourcePerm).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.resourcePerm, nil
}
