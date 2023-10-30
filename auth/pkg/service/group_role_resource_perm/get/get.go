package get

import (
	group_role_resource_perm_db "auth/pkg/database/group_role_resource_perm"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db                                  *gorm.DB
	groupRoleResourcePermResourcePermId int
}

func NewService(db *gorm.DB, groupRoleResourcePermResourcePermId int) *Service {
	return &Service{
		db:                                  db,
		groupRoleResourcePermResourcePermId: groupRoleResourcePermResourcePermId,
	}
}

func (s *Service) Get() (status int, group *group_role_resource_perm_db.GroupRoleResourcePerm, err error) {
	group = new(group_role_resource_perm_db.GroupRoleResourcePerm)
	if err = s.db.Where("GroupRoleResourcePermRoleId = ?", s.groupRoleResourcePermResourcePermId).First(group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}
	return fiber.StatusOK, group, nil
}
