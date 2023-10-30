package create

import (
	group_db "auth/pkg/database/group"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db    *gorm.DB
	group *group_db.Group
}

func NewService(db *gorm.DB, group *group_db.Group) *Service {
	return &Service{
		db:    db,
		group: group,
	}
}

func (s *Service) Create() (status int, group *group_db.Group, err error) {
	if err = s.db.Create(s.group).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fiber.StatusBadRequest, nil, err
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.group, nil
}
