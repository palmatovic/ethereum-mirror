package database

import (
	"time"
)

func (FollowedWallet) TableName() string {
	return "FollowedWallet"
}

// la clase rappresenta i muovimenti sull wallet_token/conto,
// qui sono elencate tutte le entrate ed uscite al fine di recuperare il saldo corrente
// del wallet_token/conto
// i muovimenti possono essere legati ad un aggiunta, rimozione di soldi o ad un ordine
type FollowedWallet struct {
	FollowedWalletId    int       `gorm:"column:FollowedWalletId;primaryKey;type:int;not null"`
	Balance             float64   `gorm:"column:Balance;type:float64;not null"`
	WalletIdentificator string    `gorm:"column:WalletIdentificator;varchar(255)"`
	CreatedAt           time.Time `gorm:"column:CreatedAt;type:DATETIME;autoCreateTime;not null"`
	Deleted             bool      `gorm:"column:Deleted;type:boolean;default:false"`
	DeletedAt           time.Time `gorm:"column:CreatedAt;type:DATETIME;autoCreateTime"`
}
