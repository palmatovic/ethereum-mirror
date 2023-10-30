package update

import (
	user_group_role_db "auth/pkg/database/user_group_role"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db                *gorm.DB
	userUserGroupRole *user_group_role_db.UserGroupRole
}

func NewService(db *gorm.DB, userUserGroupRole *user_group_role_db.UserGroupRole) *Service {
	return &Service{
		db:                db,
		userUserGroupRole: userUserGroupRole,
	}
}

func (s *Service) Update() (status int, group *user_group_role_db.UserGroupRole, err error) {
	if err = s.db.Where("UserGroupRoleRoleId = ?", s.userUserGroupRole.UserGroupRoleId).Updates(s.userUserGroupRole).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.userUserGroupRole, nil
}
