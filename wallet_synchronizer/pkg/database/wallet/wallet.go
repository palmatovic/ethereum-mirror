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
	Type      bool      `json:"type" gorm:"column:Type;not null"`
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

var CreateWalletGraphQL = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "create_wallet",
	Fields: graphql.InputObjectConfigFieldMap{
		"wallet_id": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
	},
})
