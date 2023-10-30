package get_by_name

import (
	product_db "auth/pkg/database/product"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db   *gorm.DB
	Name string
}

func NewService(db *gorm.DB, Name string) *Service {
	return &Service{
		db:   db,
		Name: Name,
	}
}

func (s *Service) Get() (status int, product *product_db.Product, err error) {
	product = new(product_db.Product)
	if err = s.db.Where("Name = ?", s.Name).First(product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}
	return fiber.StatusOK, product, nil
}
