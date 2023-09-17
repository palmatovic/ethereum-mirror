package token

import "time"

func (Token) TableName() string {
	return "Token"
}

type Token struct {
	TokenId        string    `json:"token_id" gorm:"column:TokenId;primaryKey"`
	Name           string    `json:"name" gorm:"column:Name;not null"`
	Symbol         string    `json:"symbol" gorm:"column:Symbol;not null"`
	Decimals       int       `json:"decimals" gorm:"column:Decimals;not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	Logo           string    `json:"logo" gorm:"column:Logo"`
	GoPlusResponse []byte    `json:"go_plus_response" gorm:"column:GoPlusResponse"`
}
