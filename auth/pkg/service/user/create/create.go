package create

import (
	user_db "auth/pkg/database/user"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db   *gorm.DB
	user *user_db.User
}

func NewService(db *gorm.DB, user *user_db.User) *Service {
	return &Service{
		db:   db,
		user: user,
	}
}

func (s *Service) Create() (status int, user *user_db.User, err error) {
	if err = s.db.Create(s.user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.user, nil
}
