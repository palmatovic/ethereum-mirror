package list

import (
	user_db "auth/pkg/database/user"
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

func (s *Service) List() (status int, users *[]user_db.User, err error) {
	users = new([]user_db.User)
	if err = s.db.Order("UserId ASC").Offset(int(s.offset)).Limit(int(s.limit)).Find(users).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, users, nil
}
