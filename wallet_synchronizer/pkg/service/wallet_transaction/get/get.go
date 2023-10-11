package get

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	wallet_transaction_db "wallet-synchronizer/pkg/database/wallet_transaction"
)

type Service struct {
	db                  *gorm.DB
	walletTransactionId string
}

func NewService(db *gorm.DB, walletTransactionId string) *Service {
	return &Service{
		db:                  db,
		walletTransactionId: walletTransactionId,
	}
}

func (s Service) Get() (status int, walletTransaction *wallet_transaction_db.WalletTransaction, err error) {
	walletTransaction = new(wallet_transaction_db.WalletTransaction)
	if err = s.db.Where("WalletTransactionId = ?", s.walletTransactionId).First(walletTransaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}
	return fiber.StatusOK, walletTransaction, nil
}
