package database

import (
	"time"
)

func (Movement) TableName() string {
	return "Movement"
}

// la clase rappresenta i muovimenti sull wallet/conto,
// qui sono elencate tutte le entrate ed uscite al fine di recuperare il saldo corrente
// del wallet/conto
// i muovimenti possono essere legati ad un aggiunta, rimozione di soldi o ad un ordine
type Movement struct {
	MovementId int     `gorm:"column:MovementId;primaryKey;type:int;not null"`
	Value      float64 `gorm:"column:Value;type:float64;not null;default:0"`
	OrderId    int     `gorm:"column:ScrapingId;type:int;null"`
	Order      Order
	CreatedAt  time.Time `gorm:"column:CreatedAt;type:DATETIME;autoCreateTime;not null"`
	Deleted    bool      `gorm:"column:Deleted;type:boolean;default:false"`
	DeletedAt  time.Time `gorm:"column:CreatedAt;type:DATETIME;autoCreateTime"`
}
