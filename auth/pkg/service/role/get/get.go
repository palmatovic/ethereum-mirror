package get

import (
	role_db "auth/pkg/database/role"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	roleId int64
}

func NewService(db *gorm.DB, roleId int64) *Service {
	return &Service{
		db:     db,
		roleId: roleId,
	}
}

func (s *Service) Get() (status int, role *role_db.Role, err error) {
	role = new(role_db.Role)
	if err = s.db.Where("RoleId = ?", s.roleId).First(role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}
	return fiber.StatusOK, role, nil
}
