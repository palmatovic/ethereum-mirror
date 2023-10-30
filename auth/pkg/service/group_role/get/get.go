package get

import (
	group_role_db "auth/pkg/database/group_role"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db          *gorm.DB
	groupRoleId int64
}

func NewService(db *gorm.DB, groupRoleId int64) *Service {
	return &Service{
		db:          db,
		groupRoleId: groupRoleId,
	}
}

func (s *Service) Get() (status int, group *group_role_db.GroupRole, err error) {
	group = new(group_role_db.GroupRole)
	if err = s.db.Where("GroupRoleRoleId = ?", s.groupRoleId).First(group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}
	return fiber.StatusOK, group, nil
}
