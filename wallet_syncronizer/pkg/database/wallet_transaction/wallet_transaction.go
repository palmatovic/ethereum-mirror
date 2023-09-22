package wallet_transaction

import (
	"time"
	"wallet-syncronizer/pkg/database/wallet"
)

func (WalletTransaction) TableName() string {
	return "WalletTransaction"
}

type WalletTransaction struct {
	TxType       string        `json:"tx_type" gorm:"column:TxType"`
	TxHash       string        `json:"tx_hash" gorm:"column:TxHash"`
	Price        float64       `json:"price" gorm:"column:Price;not null"`
	Amount       float64       `json:"amount" gorm:"column:Amount"`
	Total        float64       `json:"total" gorm:"column:Total;not null"`
	AgeTimestamp time.Time     `json:"age_timestamp" gorm:"column:AgeTimestamp;type:DATETIME;not null"`
	Asset        string        `json:"asset" gorm:"column:Asset;varchar(100)"`
	WalletId     string        `json:"wallet_id" gorm:"column:WalletId"`
	Wallet       wallet.Wallet `json:"-"`
	CreatedAt    time.Time     `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	ProcessedAt  time.Time     `json:"updated_at" gorm:"column:ProcessedAt"`
}
