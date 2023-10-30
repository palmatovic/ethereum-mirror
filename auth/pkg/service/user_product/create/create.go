package create

import (
	user_product_db "auth/pkg/database/user_product"
	"errors"
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

func (s *Service) Create() (status int, group *user_product_db.UserProduct, err error) {
	if err = s.db.Create(s.userProduct).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fiber.StatusBadRequest, nil, err
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.userProduct, nil
}
