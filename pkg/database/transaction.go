package database

import "time"

func (Transaction) TableName() string {
	return "Transaction"
}

type Transaction struct {
	TransactionHash      string `gorm:"column:TransactionHash;primaryKey;varchar(1000)"`
	ScrapingId           int    `gorm:"column:ScrapingId;type:int;not null"`
	Scraping             Scraping
	Method               string    `gorm:"column:Method;varchar(20);not null"`
	Block                string    `gorm:"column:Block;varchar(20);not null"`
	AgeMillis            string    `gorm:"column:AgeMillis;varchar(200);not null"`
	AgeTimestamp         time.Time `gorm:"column:AgeTimestamp;type:DATETIME;not null"`
	AgeDistanceFromQuery string    `gorm:"column:AgeDistanceFromQuery;varchar(200);not null"`
	GasPrice             string    `gorm:"column:GasPrice;varchar(200);not null"`
	From                 string    `gorm:"column:From;varchar(4096);not null"`
	To                   string    `gorm:"column:To;varchar(4096);not null"`
	InOut                string    `gorm:"column:InOut;varchar(200);not null"`
	Value                string    `gorm:"column:Value;varchar(200);not null"`
	TxnFee               string    `gorm:"column:TxnFee;varchar(200);not null"`
	CreatedAt            time.Time `gorm:"column:CreatedAt;type:DATETIME;autoCreateTime;not null"`
	ProcessedAt          time.Time `gorm:"column:ProcessedAt;type:DATETIME"`
}
