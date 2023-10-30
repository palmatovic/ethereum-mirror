package delete

import (
	group_role_resource_perm_db "auth/pkg/database/group_role_resource_perm"
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

func (s *Service) Delete() (status int, group *group_role_resource_perm_db.GroupRoleResourcePerm, err error) {
	group = new(group_role_resource_perm_db.GroupRoleResourcePerm)
	if err = s.db.Where("GroupRoleResourcePermId = ?", s.groupId).Delete(group).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) || errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, group, nil
}
