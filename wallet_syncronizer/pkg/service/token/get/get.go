package get

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	token_db "wallet-syncronizer/pkg/database/token"
)

type Service struct {
	db      *gorm.DB
	tokenId string
}

func NewService(db *gorm.DB, tokenId string) *Service {
	return &Service{
		db:      db,
		tokenId: tokenId,
	}
}

func (s Service) Get() (status int, token *token_db.Token, err error) {
	token = new(token_db.Token)
	if err = s.db.Where("TokenId = ?", s.tokenId).First(token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}
	return fiber.StatusOK, token, nil
}
