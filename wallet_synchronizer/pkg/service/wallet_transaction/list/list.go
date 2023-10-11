package list

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	wallet_transaction_db "wallet-synchronizer/pkg/database/wallet_transaction"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s Service) List() (status int, walletTransactions *[]wallet_transaction_db.WalletTransaction, err error) {
	walletTransactions = new([]wallet_transaction_db.WalletTransaction)
	if err = s.db.Find(walletTransactions).Error; err != nil {

		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, walletTransactions, nil
}
