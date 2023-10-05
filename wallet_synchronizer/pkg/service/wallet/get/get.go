package get

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	wallet_db "wallet-synchronizer/pkg/database/wallet"
)

type Service struct {
	db       *gorm.DB
	walletId string
}

func NewService(db *gorm.DB, walletId string) *Service {
	return &Service{
		db:       db,
		walletId: walletId,
	}
}

func (s *Service) Get() (status int, wallet *wallet_db.Wallet, err error) {
	wallet = new(wallet_db.Wallet)
	if err = s.db.Where("WalletId = ?", s.walletId).First(wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}
	return fiber.StatusOK, wallet, nil
}
