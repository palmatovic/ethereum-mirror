package list

import (
	role_db "auth/pkg/database/role"
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

func (s *Service) List() (status int, roles *[]role_db.Role, err error) {
	roles = new([]role_db.Role)
	if err = s.db.Order("Name ASC").Offset(s.offset).Limit(s.limit).Find(roles).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, roles, nil
}
