package update

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	wallet_db "wallet-synchronizer/pkg/database/wallet"
)

// Deprecated
type Service struct {
	db     *gorm.DB
	wallet *wallet_db.Wallet
}

// Deprecated
func NewService(db *gorm.DB, wallet *wallet_db.Wallet) *Service {
	return &Service{
		db:     db,
		wallet: wallet,
	}
}

// Deprecated
func (s *Service) Update() (status int, wallet *wallet_db.Wallet, err error) {
	if err = s.db.Where("WalletId = ?", s.wallet.WalletId).Updates(s.wallet).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.wallet, nil
}
