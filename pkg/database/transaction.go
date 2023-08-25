package database

import "time"

func (Transaction) TableName() string {
	return "Transaction"
}

type Transaction struct {
	TransactionHash    string `gorm:"column:TransactionHash;primaryKey;varchar(1000)"`
	ScrapingId         int    `gorm:"column:ScrapingId;type:int"`
	Scraping           Scraping
	TxType             string    `gorm:"column:TxType;varchar(20)"`
	Method             string    `gorm:"column:Method;varchar(20)"`
	Block              string    `gorm:"column:Block;varchar(20)"`
	AgeMillis          string    `gorm:"column:AgeMillis;varchar(200)"`
	AgeTimestamp       time.Time `gorm:"column:AgeTimestamp;type:DATETIME"`
	AgeDistanceFromNow string    `gorm:"column:AgeDistanceFromNow;varchar(200)"`
	GasPrice           string    `gorm:"column:GasPrice;varchar(200)"`
	From               string    `gorm:"column:From;varchar(200)"`
	To                 string    `gorm:"column:To;varchar(200)"`
	Amount             string    `gorm:"column:Amount;varchar(200)"`
	Value              string    `gorm:"column:Value;varchar(200)"`
	Asset              string    `gorm:"column:Asset;varchar(200)"`
	TxnFee             string    `gorm:"column:TxnFee;varchar(200)"`
	CreatedAt          time.Time `gorm:"column:CreatedAt;type:DATETIME;autoCreateTime"`
	ProcessedAt        time.Time `gorm:"column:ProcessedAt;type:DATETIME"`
}
