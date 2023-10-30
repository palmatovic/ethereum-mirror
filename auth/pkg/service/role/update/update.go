package update

import (
	role_db "auth/pkg/database/role"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db   *gorm.DB
	role *role_db.Role
}

func NewService(db *gorm.DB, role *role_db.Role) *Service {
	return &Service{
		db:   db,
		role: role,
	}
}

func (s *Service) Update() (status int, role *role_db.Role, err error) {
	if err = s.db.Where("RoleId = ?", s.role.RoleId).Updates(s.role).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.role, nil
}
