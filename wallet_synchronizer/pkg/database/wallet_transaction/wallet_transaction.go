package wallet_transaction

import (
	"github.com/graphql-go/graphql"
	"time"
	"wallet-synchronizer/pkg/database/wallet"
)

func (WalletTransaction) TableName() string {
	return "WalletTransaction"
}

type WalletTransaction struct {
	WalletTransactionId string        `json:"tx_hash" gorm:"column:WalletTransactionId;primaryKey"`
	TxType              string        `json:"tx_type" gorm:"column:TxType"`
	Price               float64       `json:"price" gorm:"column:Price;not null"`
	Amount              float64       `json:"amount" gorm:"column:Amount"`
	Total               float64       `json:"total" gorm:"column:Total;not null"`
	AgeTimestamp        time.Time     `json:"age_timestamp" gorm:"column:AgeTimestamp;type:DATETIME;not null"`
	Asset               string        `json:"asset" gorm:"column:Asset;varchar(100)"`
	WalletId            string        `json:"wallet_id" gorm:"column:WalletId"`
	Wallet              wallet.Wallet `json:"-"`
	CreatedAt           time.Time     `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
}

var WalletTransactionGraphQL = graphql.NewObject(graphql.ObjectConfig{
	Name: "wallet_transaction",
	Fields: graphql.Fields{
		"wallet_transaction_id": &graphql.Field{
			Type: graphql.String,
		},
		"tx_type": &graphql.Field{
			Type: graphql.DateTime,
		},
		"price": &graphql.Field{
			Type: graphql.Float,
		},
		"amount": &graphql.Field{
			Type: graphql.Float,
		},
		"total": &graphql.Field{
			Type: graphql.Float,
		},
		"age_timestamp": &graphql.Field{
			Type: graphql.DateTime,
		},

		"asset": &graphql.Field{
			Type: graphql.String,
		},
		"wallet_id": &graphql.Field{
			Type: graphql.String,
		},
		"created_at": &graphql.Field{
			Type: graphql.DateTime,
		},
	},
})
