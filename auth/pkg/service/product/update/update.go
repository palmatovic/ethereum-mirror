package update

import (
	product_db "auth/pkg/database/product"
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

func (s *Service) Update() (status int, product *product_db.Product, err error) {
	if err = s.db.Where("ProductId = ?", s.product.ProductId).Updates(s.product).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.product, nil
}
