package update

import (
	group_role_db "auth/pkg/database/group_role"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	groupRole *group_role_db.GroupRole
}

func NewService(db *gorm.DB, groupRole *group_role_db.GroupRole) *Service {
	return &Service{
		db:        db,
		groupRole: groupRole,
	}
}

func (s *Service) Update() (status int, group *group_role_db.GroupRole, err error) {
	if err = s.db.Where("GroupRoleRoleId = ?", s.groupRole.GroupRoleId).Updates(s.groupRole).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.groupRole, nil
}
