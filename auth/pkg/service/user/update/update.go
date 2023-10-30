package update

import (
	user_db "auth/pkg/database/user"
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

func (s *Service) Update() (status int, user *user_db.User, err error) {
	if err = s.db.Where("UserId = ?", s.user.UserId).Updates(s.user).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.user, nil
}
