package list

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	wallet_db "wallet-synchronizer/pkg/database/wallet"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) List() (status int, wallets *[]wallet_db.Wallet, err error) {
	wallets = new([]wallet_db.Wallet)
	if err = s.db.Find(wallets).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, wallets, nil
}
