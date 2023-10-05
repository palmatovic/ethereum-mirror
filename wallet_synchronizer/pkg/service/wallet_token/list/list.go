package list

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	wallet_token_db "wallet-synchronizer/pkg/database/wallet_token"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s Service) List() (status int, walletTokens *[]wallet_token_db.WalletToken, err error) {
	walletTokens = new([]wallet_token_db.WalletToken)
	if err = s.db.Find(walletTokens).Error; err != nil {

		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, walletTokens, nil
}
