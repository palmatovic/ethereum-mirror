package update

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	token_db "wallet-synchronizer/pkg/database/token"
)

type Service struct {
	db    *gorm.DB
	token *token_db.Token
}

func NewService(db *gorm.DB, token *token_db.Token) *Service {
	return &Service{
		db:    db,
		token: token,
	}
}

func (s Service) Update() (status int, token *token_db.Token, err error) {
	if err = s.db.Where("TokenId = ?", s.token.TokenId).Updates(s.token).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.token, nil
}
