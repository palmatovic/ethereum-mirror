package database

import "time"

func (Scraping) TableName() string {
	return "Scraping"
}

type Scraping struct {
	ScrapingId    int       `gorm:"column:ScrapingId;primaryKey;autoIncrement"`
	WalletAddress string    `gorm:"column:WalletAddress;varchar(1024)"`
	CreatedAt     time.Time `gorm:"column:AgeTimestamp;type:DATETIME;autoCreateTime"`
}
