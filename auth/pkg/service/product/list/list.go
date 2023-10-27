package list

import (
	product_db "auth/pkg/database/product"
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

func (s *Service) List() (status int, products *[]product_db.Product, err error) {
	products = new([]product_db.Product)
	if err = s.db.Order("Name ASC").Offset(s.offset).Limit(s.limit).Find(products).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, products, nil
}
