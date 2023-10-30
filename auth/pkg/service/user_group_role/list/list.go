package list

import (
	user_group_role_db "auth/pkg/database/user_group_role"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	limit  int
	offset int
}

func NewService(db *gorm.DB, pageSize int, pageNumber int) *Service {
	return &Service{
		db:     db,
		limit:  pageSize,
		offset: (pageNumber - 1) * pageSize,
	}
}

func (s *Service) List() (status int, groups *[]user_group_role_db.UserGroupRole, err error) {
	groups = new([]user_group_role_db.UserGroupRole)
	if err = s.db.Order("UserGroupRoleId ASC").Offset(s.offset).Limit(s.limit).Find(groups).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, groups, nil
}
