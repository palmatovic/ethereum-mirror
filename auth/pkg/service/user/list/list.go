package list

import (
	user_db "auth/pkg/database/user"
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

func (s *Service) List() (status int, users *[]user_db.User, err error) {
	users = new([]user_db.User)
	if err = s.db.Order("UserId ASC").Offset(s.offset).Limit(s.limit).Find(users).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, users, nil
}
