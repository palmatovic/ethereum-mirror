package create

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	wallet_token_db "wallet-syncronizer/pkg/database/wallet_token"
)

type Service struct {
	db          *gorm.DB
	walletToken *wallet_token_db.WalletToken
}

func NewService(db *gorm.DB, walletToken *wallet_token_db.WalletToken) *Service {
	return &Service{
		db:          db,
		walletToken: walletToken,
	}
}

func (s Service) Create() (status int, token *wallet_token_db.WalletToken, err error) {
	if err = s.db.Create(s.walletToken).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.walletToken, nil
}
