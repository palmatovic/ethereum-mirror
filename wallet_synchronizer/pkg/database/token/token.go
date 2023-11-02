package token

import (
	"encoding/json"
	"github.com/graphql-go/graphql"
	"time"
)

type Token struct {
	TokenId        string    `json:"token_id" gorm:"column:TokenId;primaryKey"`
	Name           string    `json:"name" gorm:"column:Name;not null"`
	Symbol         string    `json:"symbol" gorm:"column:Symbol;not null"`
	Decimals       int       `json:"decimals" gorm:"column:Decimals;not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	Logo           string    `json:"logo" gorm:"column:Logo"`
	CurrentPrice   float64   `json:"current_price" gorm:"column:CurrentPrice"`
	GoPlusResponse []byte    `json:"go_plus_response" gorm:"column:GoPlusResponse"`
}

func (Token) TableName() string {
	return "Token"
}

var TokenGraphQL = graphql.NewObject(graphql.ObjectConfig{
	Name: "token",
	Fields: graphql.Fields{
		"token_id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"symbol": &graphql.Field{
			Type: graphql.String,
		},
		"decimals": &graphql.Field{
			Type: graphql.Int,
		},
		"created_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"logo": &graphql.Field{
			Type: graphql.String,
		},
		"current_price": &graphql.Field{
			Type: graphql.Float,
		},
		"go_plus_response": &graphql.Field{
			Type: graphql.NewScalar(graphql.ScalarConfig{
				Name: "Json",
				Serialize: func(value interface{}) interface{} {
					var serialized map[string]interface{}
					_ = json.Unmarshal(value.([]byte), &serialized)
					return serialized
				},
			}),
		},
	},
})
