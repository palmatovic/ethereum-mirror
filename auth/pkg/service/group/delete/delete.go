package delete

import (
	group_db "auth/pkg/database/group"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db      *gorm.DB
	groupId int64
}

func NewService(db *gorm.DB, groupId int64) *Service {
	return &Service{
		db:      db,
		groupId: groupId,
	}
}

func (s *Service) Delete() (status int, group *group_db.Group, err error) {
	group = new(group_db.Group)
	if err = s.db.Where("GroupId = ?", s.groupId).Delete(group).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) || errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, group, nil
}
