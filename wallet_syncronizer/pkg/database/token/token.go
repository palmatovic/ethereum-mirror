package token

import "time"

func (Token) TableName() string {
	return "Token"
}

type Token struct {
	TokenId     string    `gorm:"column:TokenId;primaryKey"`
	Name        string    `gorm:"column:Name;not null"`
	Symbol      string    `gorm:"column:Symbol;not null"`
	Decimals    int       `gorm:"column:Decimals;not null"`
	CreatedAt   time.Time `gorm:"column:CreatedAt;autoCreateTime"`
	Logo        string    `gorm:"column:Logo"`
	RiskScam    int       `gorm:"column:RiskScam"`
	WarningScam int       `gorm:"column:WarningScam"`
}
