package find_or_create

import (
	"errors"
	"gorm.io/gorm"
	wallet_db "wallet-synchronizer/pkg/database/wallet"
	wallet_create_service "wallet-synchronizer/pkg/service/wallet/create"
	wallet_get_service "wallet-synchronizer/pkg/service/wallet/get"
)

type Service struct {
	db            *gorm.DB
	walletAddress string
}

func NewService(db *gorm.DB, walletAddress string) *Service {
	return &Service{
		db:            db,
		walletAddress: walletAddress,
	}
}

func (s *Service) FindOrCreateWalletToMonitor() (wallet *wallet_db.Wallet, err error) {
	if _, wallet, err = wallet_get_service.NewService(s.db, s.walletAddress).Get(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if _, wallet, err = wallet_create_service.NewService(s.db, &wallet_db.Wallet{WalletId: s.walletAddress, Type: true}).Create(); err != nil {
				return wallet, err
			}
		}
	}
	return wallet, nil
}
