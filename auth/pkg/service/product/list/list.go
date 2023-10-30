package list

import (
	product_db "auth/pkg/database/product"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	limit  int64
	offset int64
}

func NewService(db *gorm.DB, pageSize int64, pageNumber int64) *Service {
	return &Service{
		db:     db,
		limit:  pageSize,
		offset: (pageNumber - 1) * pageSize,
	}
}

func (s *Service) List() (status int, products *[]product_db.Product, err error) {
	products = new([]product_db.Product)
	if err = s.db.Order("Name ASC").Offset(int(s.offset)).Limit(int(s.limit)).Find(products).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, products, nil
}
