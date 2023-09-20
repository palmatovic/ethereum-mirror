package wallet

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	wallet_db "wallet-syncronizer/pkg/database/wallet"
)

func FindOrCreateWallet(walletAddress string, db *gorm.DB) (wallet wallet_db.Wallet, err error) {
	if err = db.Where("WalletId = ?", walletAddress).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			wallet = wallet_db.Wallet{
				WalletId: walletAddress,
			}
			if err = db.Create(&wallet).Error; err != nil {
				return wallet, err
			}
		}
	}
	return wallet, nil
}

func GetWallet(db *gorm.DB, walletId string) (status int, wallet *wallet_db.Wallet, err error) {
	wallet = new(wallet_db.Wallet)
	if err = db.Where("WalletId = ?", walletId).First(wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}
	return fiber.StatusOK, wallet, nil
}
