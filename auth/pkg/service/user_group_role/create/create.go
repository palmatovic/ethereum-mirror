package create

import (
	user_group_role_db "auth/pkg/database/user_group_role"
	"errors"
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

func (s *Service) Create() (status int, group *user_group_role_db.UserGroupRole, err error) {
	if err = s.db.Create(s.userUserGroupRole).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fiber.StatusBadRequest, nil, err
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.userUserGroupRole, nil
}
