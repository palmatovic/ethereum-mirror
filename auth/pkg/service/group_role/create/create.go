package create

import (
	group_role_db "auth/pkg/database/group_role"
	"errors"
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

func (s *Service) Create() (status int, group *group_role_db.GroupRole, err error) {
	if err = s.db.Create(s.groupRole).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fiber.StatusBadRequest, nil, err
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.groupRole, nil
}
