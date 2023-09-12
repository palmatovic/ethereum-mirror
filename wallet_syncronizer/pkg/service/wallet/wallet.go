package wallet

import (
	"errors"
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
