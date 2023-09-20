package token

import (
	"reflect"
	"strings"
	"time"
)

type Token struct {
	TokenId        string    `json:"token_id" gorm:"column:TokenId;primaryKey"`
	Name           string    `json:"name" gorm:"column:Name;not null"`
	Symbol         string    `json:"symbol" gorm:"column:Symbol;not null"`
	Decimals       int       `json:"decimals" gorm:"column:Decimals;not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	Logo           string    `json:"logo" gorm:"column:Logo"`
	GoPlusResponse []byte    `json:"go_plus_response" gorm:"column:GoPlusResponse"`
}

func (Token) TableName() string {
	return "Token"
}

func (t *Token) GetColumnName(field interface{}) string {
	v := reflect.ValueOf(t).Elem()
	for i := 0; i < v.NumField(); i++ {
		if reflect.DeepEqual(v.Field(i).Interface(), reflect.ValueOf(field).Elem().Interface()) {
			tag := v.Type().Field(i).Tag.Get("gorm")
			parts := strings.Split(tag, ";")
			for _, part := range parts {
				if strings.HasPrefix(part, "column:") {
					return strings.TrimPrefix(part, "column:")
				}
			}
		}
	}
	return ""
}
