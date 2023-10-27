package create

import (
	product_db "auth/pkg/database/product"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db      *gorm.DB
	product *product_db.Product
}

func NewService(db *gorm.DB, product *product_db.Product) *Service {
	return &Service{
		db:      db,
		product: product,
	}
}

func (s *Service) Create() (status int, product *product_db.Product, err error) {
	if err = s.db.Create(s.product).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.product, nil
}
