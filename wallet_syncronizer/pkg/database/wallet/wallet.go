package wallet

import "time"

func (Wallet) TableName() string {
	return "Wallet"
}

type Wallet struct {
	WalletId  string    `gorm:"column:WalletId;primaryKey;varchar(1024)"`
	CreatedAt time.Time `gorm:"column:CreatedAt"`
}
