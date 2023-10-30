package list

import (
	group_db "auth/pkg/database/group"
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

func (s *Service) List() (status int, groups *[]group_db.Group, err error) {
	groups = new([]group_db.Group)
	if err = s.db.Order("Name ASC").Offset(s.offset).Limit(s.limit).Find(groups).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, groups, nil
}
