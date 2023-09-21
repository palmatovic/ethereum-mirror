package list

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	token_db "wallet-syncronizer/pkg/database/token"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s Service) List() (status int, tokens *[]token_db.Token, err error) {
	tokens = new([]token_db.Token)
	if err = s.db.Find(tokens).Error; err != nil {

		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, tokens, nil
}
