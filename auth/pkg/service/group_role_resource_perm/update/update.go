package update

import (
	group_role_resource_perm_db "auth/pkg/database/group_role_resource_perm"
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

func (s *Service) Update() (status int, group *group_role_resource_perm_db.GroupRoleResourcePerm, err error) {
	if err = s.db.Where("GroupRoleResourcePermRoleId = ?", s.groupRoleResourcePermResourcePerm.GroupRoleResourcePermId).Updates(s.groupRoleResourcePermResourcePerm).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.groupRoleResourcePermResourcePerm, nil
}
