package delete

import (
	resource_perm_db "auth/pkg/database/resource_perm"
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

func (s *Service) Delete() (status int, group *resource_perm_db.ResourcePerm, err error) {
	group = new(resource_perm_db.ResourcePerm)
	if err = s.db.Where("ResourcePermId = ?", s.groupId).Delete(group).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) || errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, group, nil
}
