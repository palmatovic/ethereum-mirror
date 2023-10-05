package update

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	wallet_token_db "wallet-synchronizer/pkg/database/wallet_token"
)

type Service struct {
	db          *gorm.DB
	walletToken *wallet_token_db.WalletToken
}

func NewService(db *gorm.DB, walletToken *wallet_token_db.WalletToken) *Service {
	return &Service{
		db:          db,
		walletToken: walletToken,
	}
}

func (s Service) Update() (status int, walletToken *wallet_token_db.WalletToken, err error) {
	if err = s.db.Where("WalletId = ? AND TokenId", s.walletToken.WalletId, s.walletToken.TokenId).Updates(s.walletToken).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.walletToken, nil
}
