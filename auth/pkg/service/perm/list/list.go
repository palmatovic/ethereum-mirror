package list

import (
	perm_db "auth/pkg/database/perm"
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

func (s *Service) List() (status int, perms *[]perm_db.Perm, err error) {
	perms = new([]perm_db.Perm)
	if err = s.db.Order("PermId ASC").Offset(int(s.offset)).Limit(int(s.limit)).Find(perms).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, perms, nil
}
