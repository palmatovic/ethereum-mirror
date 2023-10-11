package wallet

import (
	"github.com/graphql-go/graphql"
	"time"
)

func (Wallet) TableName() string {
	return "Wallet"
}

type Wallet struct {
	WalletId  string    `json:"wallet_id" gorm:"column:WalletId;primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
}

var WalletGraphQL = graphql.NewObject(graphql.ObjectConfig{
	Name: "wallet",
	Fields: graphql.Fields{
		"wallet_id": &graphql.Field{
			Type: graphql.String,
		},
		"created_at": &graphql.Field{
			Type: graphql.DateTime,
		},
	},
})
