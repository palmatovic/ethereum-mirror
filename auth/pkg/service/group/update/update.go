package update

import (
	group_db "auth/pkg/database/group"
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

func (s *Service) Update() (status int, group *group_db.Group, err error) {
	if err = s.db.Where("GroupId = ?", s.group.GroupId).Updates(s.group).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.group, nil
}
