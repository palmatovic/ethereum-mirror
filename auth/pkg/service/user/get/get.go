package get

import (
	user_db "auth/pkg/database/user"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	userId int64
}

func NewService(db *gorm.DB, userId int64) *Service {
	return &Service{
		db:     db,
		userId: userId,
	}
}

func (s *Service) Get() (status int, user *user_db.User, err error) {
	user = new(user_db.User)
	if err = s.db.Where("UserId = ?", s.userId).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}
	return fiber.StatusOK, user, nil
}
