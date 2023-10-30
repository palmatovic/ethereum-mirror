package create

import (
	group_role_resource_perm_db "auth/pkg/database/group_role_resource_perm"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db                                *gorm.DB
	groupRoleResourcePermResourcePerm *group_role_resource_perm_db.GroupRoleResourcePerm
}

func NewService(db *gorm.DB, groupRoleResourcePermResourcePerm *group_role_resource_perm_db.GroupRoleResourcePerm) *Service {
	return &Service{
		db:                                db,
		groupRoleResourcePermResourcePerm: groupRoleResourcePermResourcePerm,
	}
}

func (s *Service) Create() (status int, group *group_role_resource_perm_db.GroupRoleResourcePerm, err error) {
	if err = s.db.Create(s.groupRoleResourcePermResourcePerm).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fiber.StatusBadRequest, nil, err
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.groupRoleResourcePermResourcePerm, nil
}
