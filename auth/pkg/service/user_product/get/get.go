package get

import (
	user_product_db "auth/pkg/database/user_product"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db            *gorm.DB
	userProductId int
}

func NewService(db *gorm.DB, userProductId int) *Service {
	return &Service{
		db:            db,
		userProductId: userProductId,
	}
}

func (s *Service) Get() (status int, group *user_product_db.UserProduct, err error) {
	group = new(user_product_db.UserProduct)
	if err = s.db.Where("UserProductRoleId = ?", s.userProductId).First(group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}
	return fiber.StatusOK, group, nil
}
