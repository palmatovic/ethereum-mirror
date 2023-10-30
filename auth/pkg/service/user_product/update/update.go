package update

import (
	user_product_db "auth/pkg/database/user_product"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db          *gorm.DB
	userProduct *user_product_db.UserProduct
}

func NewService(db *gorm.DB, userProduct *user_product_db.UserProduct) *Service {
	return &Service{
		db:          db,
		userProduct: userProduct,
	}
}

func (s *Service) Update() (status int, group *user_product_db.UserProduct, err error) {
	if err = s.db.Where("UserProductRoleId = ?", s.userProduct.UserProductId).Updates(s.userProduct).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.userProduct, nil
}
