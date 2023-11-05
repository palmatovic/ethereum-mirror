package token

import (
	"time"

	"github.com/graphql-go/graphql"
)

func (TokenPrice) TableName() string {
	return "TokenPrice"
}

type TokenPrice struct {
	TokenPriceId int       `json:"token_price_id" gorm:"column:TokenPriceId;primaryKey;type:int;not null"`
	TokenId      string    `json:"token_id" gorm:"column:TokenId;primaryKey;not null"`
	Token        Token     `json:"-"`
	Price        float64   `json:"price" gorm:"column:price;not null"`
	PriceDate    time.Time `json:"price_date" gorm:"column:PriceDate;not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
}

var TokenPriceGraphQL = graphql.NewObject(graphql.ObjectConfig{
	Name: "token_price",
	Fields: graphql.Fields{
		"token_price_id": &graphql.Field{
			Type: graphql.Int,
		},
		"token_id": &graphql.Field{
			Type: graphql.String,
		},
		"price": &graphql.Field{
			Type: graphql.Float,
		},
		"price_date": &graphql.Field{
			Type: graphql.DateTime,
		},
		"created_at": &graphql.Field{
			Type: graphql.DateTime,
		},
	},
})
