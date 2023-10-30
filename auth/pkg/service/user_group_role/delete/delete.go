package delete

import (
	user_group_role_db "auth/pkg/database/user_group_role"
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

func (s *Service) Delete() (status int, group *user_group_role_db.UserGroupRole, err error) {
	group = new(user_group_role_db.UserGroupRole)
	if err = s.db.Where("UserGroupRoleId = ?", s.groupId).Delete(group).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) || errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, group, nil
}
