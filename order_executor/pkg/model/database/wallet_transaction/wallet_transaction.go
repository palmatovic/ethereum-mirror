package wallet_transaction

import (
	"order-executor/pkg/model/database/token"
	"order-executor/pkg/model/database/wallet"
	"time"

	"github.com/graphql-go/graphql"
)

func (WalletTransaction) TableName() string {
	return "WalletTransaction"
}

type WalletTransaction struct {
	WalletTransactionId            string    `json:"tx_hash" gorm:"column:WalletTransactionId;primaryKey"`
	TxType                         string    `json:"tx_type" gorm:"column:TxType"`
	Pool                           string    `json:"pool" gorm:"column:Pool"`
	Price                          float64   `json:"price" gorm:"column:Price;not null"`
	Amount                         float64   `json:"amount" gorm:"column:Amount"`
	Total                          float64   `json:"total" gorm:"column:Total;not null"`
	AgeTimestamp                   time.Time `json:"age_timestamp" gorm:"column:AgeTimestamp;type:DATETIME;not null"`
	TokenId                        string    `json:"token_id" gorm:"column:TokenId"`
	Token                          token.Token
	WalletId                       string        `json:"wallet_id" gorm:"column:WalletId"`
	Wallet                         wallet.Wallet `json:"-"`
	CreatedAt                      time.Time     `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt                      time.Time     `json:"updated_at" gorm:"column:UpdatedAt;autoCreateTime"`
	ProcessedByOrderExecutor       bool          `json:"processed_by_order_executor" gorm:"column:ProcessedByOrderExecutor;default:0"`
	ProcessedByOrderExecutorAt     time.Time     `json:"processed_by_order_executor_at" gorm:"column:ProcessedByOrderExecutorAt;type:DATETIME"`
	ProcessedByOrderExecutorStatus string        `json:"processed_by_order_executor_status" gorm:"column:ProcessedByOrderExecutorStatus;default:registered"`
}

var WalletTransactionGraphQL = graphql.NewObject(graphql.ObjectConfig{
	Name: "wallet_transaction",
	Fields: graphql.Fields{
		"wallet_transaction_id": &graphql.Field{
			Type: graphql.String,
		},
		"tx_type": &graphql.Field{
			Type: graphql.String,
		},
		"pool": &graphql.Field{
			Type: graphql.String,
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

		"token_id": &graphql.Field{
			Type: graphql.String,
		},
		"wallet_id": &graphql.Field{
			Type: graphql.String,
		},
		"created_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"updated_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"processed_by_order_executor": &graphql.Field{
			Type: graphql.Boolean,
		},
		"processed_by_order_executor_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"processed_by_order_executor_status": &graphql.Field{
			Type: graphql.String,
		},
	},
})
