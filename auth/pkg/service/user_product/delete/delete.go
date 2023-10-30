package delete

import (
	user_product_db "auth/pkg/database/user_product"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db      *gorm.DB
	groupId int
}

func NewService(db *gorm.DB, groupId int) *Service {
	return &Service{
		db:      db,
		groupId: groupId,
	}
}

func (s *Service) Delete() (status int, group *user_product_db.UserProduct, err error) {
	group = new(user_product_db.UserProduct)
	if err = s.db.Where("UserProductId = ?", s.groupId).Delete(group).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) || errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, group, nil
}
