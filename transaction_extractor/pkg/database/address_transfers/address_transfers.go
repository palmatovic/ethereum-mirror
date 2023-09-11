package address_transfers

import "time"

func (Transaction) TableName() string {
	return "Transaction"
}

type Transaction struct {
	TxType       string    `gorm:"column:TxType;varchar(100)"`
	Price        float64   `gorm:"column:Price;not null"`
	Amount       float64   `gorm:"column:Amount"`
	Total        float64   `gorm:"column:Total;not null"`
	AgeTimestamp time.Time `gorm:"column:AgeTimestamp;type:DATETIME;not null"`
	//AgeDistanceFromQuery string    `gorm:"column:AgeDistanceFromQuery;varchar(200);not null"`
	//AgeMillis            string    `gorm:"column:AgeMillis;varchar(200);not null"`
	//From                 string    `gorm:"column:From;varchar(4096);not null"`
	//To                   string    `gorm:"column:To;varchar(4096);not null"`
	//Value                string    `gorm:"column:Value;varchar(200);not null"`
	Asset string `gorm:"column:Asset;varchar(100)"`
	//TxnFee   string `gorm:"column:TxnFee;varchar(200);not null"`
	//GasPrice string `gorm:"column:GasPrice;varchar(200);not null"`
	// ScrapingId                 int       `gorm:"column:ScrapingId;type:int;not null"`
	// Scraping                   scraping.Scraping
	WalletAddress string    `gorm:"column:WalletAddress;varchar(1024);not null"`
	CreatedAt     time.Time `gorm:"column:CreatedAt;type:DATETIME;autoCreateTime;not null"`
	ProcessedAt   time.Time `gorm:"column:ProcessedAt;type:DATETIME"`
	//ProcessedByOrderExecutor   bool      `gorm:"column:ProcessedByOrderExecutor;type:boolean;default:false"`
	//ProcessedByOrderExecutorAt time.Time `gorm:"column:ProcessedByOrderExecutorAt;type:DATETIME"`
	//FollowedWalletId           int       `gorm:"column:FollowedWalletId;type:int;not null"`
}
