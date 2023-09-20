package wallet

import "time"

func (Wallet) TableName() string {
	return "Wallet"
}

type Wallet struct {
	WalletId  string    `json:"wallet_id" gorm:"column:WalletId;primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
}
