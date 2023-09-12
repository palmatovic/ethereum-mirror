package wallet_token

import (
	"time"
	token_db "wallet-syncronizer/pkg/database/token"
	"wallet-syncronizer/pkg/database/wallet"
)

func (WalletToken) TableName() string {
	return "WalletToken"
}

type WalletToken struct {
	WalletId       string `gorm:"column:WalletId;primaryKey;varchar(1024)"`
	Wallet         wallet.Wallet
	TokenId        string `gorm:"column:TokenId;primaryKey;varchar(1024);not null"`
	Token          token_db.Token
	TokenAmount    float64   `gorm:"column:TokenAmount;not null"`
	TokenAmountHex string    `gorm:"column:TokenAmountHex;varchar(1024);not null"`
	CreatedAt      time.Time `gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:UpdatedAt;autoUpdateTime"`
}
