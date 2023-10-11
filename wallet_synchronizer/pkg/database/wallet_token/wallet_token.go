package wallet_token

import (
	"github.com/graphql-go/graphql"
	"time"
	token_db "wallet-synchronizer/pkg/database/token"
	"wallet-synchronizer/pkg/database/wallet"
)

func (WalletToken) TableName() string {
	return "WalletToken"
}

type WalletToken struct {
	WalletId       string         `json:"wallet_id" gorm:"column:WalletId;primaryKey"`
	Wallet         wallet.Wallet  `json:"-"`
	TokenId        string         `json:"token_id" gorm:"column:TokenId;primaryKey;not null"`
	Token          token_db.Token `json:"-"`
	TokenAmount    float64        `json:"token_amount" gorm:"column:TokenAmount;not null"`
	TokenAmountHex string         `json:"token_amount_hex" gorm:"column:TokenAmountHex;not null"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

var WalletTokenGraphQL = graphql.NewObject(graphql.ObjectConfig{
	Name: "wallet_token",
	Fields: graphql.Fields{
		"wallet_id": &graphql.Field{
			Type: graphql.String,
		},
		"token_id": &graphql.Field{
			Type: graphql.String,
		},
		"token_amount": &graphql.Field{
			Type: graphql.DateTime,
		},

		"token_amount_hex": &graphql.Field{
			Type: graphql.DateTime,
		},
		"created_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"updated_at": &graphql.Field{
			Type: graphql.DateTime,
		},
	},
})
