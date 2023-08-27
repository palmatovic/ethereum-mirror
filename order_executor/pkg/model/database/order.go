package database

import (
	"time"
)

func (Order) TableName() string {
	return "Order"
}

// la clase rappresenta gli ordini di acquisto con il relativo stato,
// al fine di monitorare le posizioni aperte ed il loro stato
type Order struct {
	OrderId       int         `gorm:"column:OrderId;primaryKey;type:int;not null"`
	Value         float64     `gorm:"column:Value;type:float64;not null"`
	CreatedAt     time.Time   `gorm:"column:CreatedAt;type:DATETIME;autoCreateTime;not null"`
	ExecutedAt    time.Time   `gorm:"column:ExecutedAt;type:DATETIME;"`
	ClosedAt      time.Time   `gorm:"column:ExecutedAt;type:DATETIME;"`
	Status        int         `gorm:"column:Status;type:DATETIME;"`                     // 1 creato, 2 inserito, 3 chiuso
	TransactionId string      `gorm:"column:TransactionId;type:varchar(1000);not null"` // This corresponds to the TxHash of the Transaction struct
	Transaction   Transaction `gorm:"foreignKey:TransactionId;references:TxHash"`       // This sets up the foreign key relationship
	Deleted       bool        `gorm:"column:Deleted;type:boolean;default:false"`
	DeletedAt     time.Time   `gorm:"column:CreatedAt;type:DATETIME;autoCreateTime"`
}
