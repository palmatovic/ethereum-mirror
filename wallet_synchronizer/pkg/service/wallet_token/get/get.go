package get

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	wallet_token_db "wallet-synchronizer/pkg/database/wallet_token"
)

type Service struct {
	db       *gorm.DB
	walletId string
	tokenId  string
}

func NewService(db *gorm.DB, walletId, tokenId string) *Service {
	return &Service{
		db:       db,
		walletId: walletId,
		tokenId:  tokenId,
	}
}

func (s Service) Get() (status int, walletToken *wallet_token_db.WalletToken, err error) {
	walletToken = new(wallet_token_db.WalletToken)
	if err = s.db.Where("WalletId = ? AND TokenId = ?", s.walletId, s.tokenId).First(walletToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}
	return fiber.StatusOK, walletToken, nil
}
