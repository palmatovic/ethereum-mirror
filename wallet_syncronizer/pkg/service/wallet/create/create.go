package create

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	wallet_db "wallet-syncronizer/pkg/database/wallet"
)

type Service struct {
	db     *gorm.DB
	wallet *wallet_db.Wallet
}

func NewService(db *gorm.DB, wallet *wallet_db.Wallet) *Service {
	return &Service{
		db:     db,
		wallet: wallet,
	}
}

func (s *Service) Create() (status int, wallet *wallet_db.Wallet, err error) {
	if err = s.db.Create(s.wallet).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.wallet, nil
}
