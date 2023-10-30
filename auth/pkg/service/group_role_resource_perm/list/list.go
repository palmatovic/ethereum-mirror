package list

import (
	group_role_resource_perm_db "auth/pkg/database/group_role_resource_perm"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	limit  int64
	offset int64
}

func NewService(db *gorm.DB, pageSize int64, pageNumber int64) *Service {
	return &Service{
		db:     db,
		limit:  pageSize,
		offset: (pageNumber - 1) * pageSize,
	}
}

func (s *Service) List() (status int, groups *[]group_role_resource_perm_db.GroupRoleResourcePerm, err error) {
	groups = new([]group_role_resource_perm_db.GroupRoleResourcePerm)
	if err = s.db.Order("GroupRoleResourcePermId ASC").Offset(int(s.offset)).Limit(int(s.limit)).Find(groups).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, groups, nil
}
