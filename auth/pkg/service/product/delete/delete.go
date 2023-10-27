package delete

import (
	product_db "auth/pkg/database/product"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	productId int
}

func NewService(db *gorm.DB, productId int) *Service {
	return &Service{
		db:        db,
		productId: productId,
	}
}

func (s *Service) Delete() (status int, product *product_db.Product, err error) {
	product = new(product_db.Product)
	if err = s.db.Where("ProductId = ?", s.productId).Delete(product).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) || errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusBadRequest, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, product, nil
}
