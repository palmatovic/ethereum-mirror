package get

import (
	group_db "auth/pkg/database/group"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db      *gorm.DB
	groupId int
}

func NewService(db *gorm.DB, groupId int) *Service {
	return &Service{
		db:      db,
		groupId: groupId,
	}
}

func (s *Service) Get() (status int, group *group_db.Group, err error) {
	group = new(group_db.Group)
	if err = s.db.Where("GroupId = ?", s.groupId).First(group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}
	return fiber.StatusOK, group, nil
}
