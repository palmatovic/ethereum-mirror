package get

import (
	user_group_role_db "auth/pkg/database/user_group_role"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db                  *gorm.DB
	userUserGroupRoleId int64
}

func NewService(db *gorm.DB, userUserGroupRoleId int64) *Service {
	return &Service{
		db:                  db,
		userUserGroupRoleId: userUserGroupRoleId,
	}
}

func (s *Service) Get() (status int, group *user_group_role_db.UserGroupRole, err error) {
	group = new(user_group_role_db.UserGroupRole)
	if err = s.db.Where("UserGroupRoleRoleId = ?", s.userUserGroupRoleId).First(group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}
	return fiber.StatusOK, group, nil
}
