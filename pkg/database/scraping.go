package database

import "time"

func (Scraping) TableName() string {
	return "Scraping"
}

type Scraping struct {
	ScrapingId  int       `gorm:"column:ScrapingId;primaryKey;autoIncrement"`
	CreatedAt   time.Time `gorm:"column:AgeTimestamp;type:DATETIME;autoCreateTime"`
	ProcessedAt time.Time `gorm:"column:ProcessedAt;type:DATETIME"`
}
