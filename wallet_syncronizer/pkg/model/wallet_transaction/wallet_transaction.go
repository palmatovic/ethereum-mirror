package wallet_transaction

import "time"

type AddressTransaction struct {
	TxType        string    `json:"tx_type"`
	Price         float64   `json:"price"`
	Amount        float64   `json:"amount"`
	Total         float64   `json:"total"`
	AgeTimestamp  time.Time `json:"age_timestamp"`
	Asset         string    `json:"asset"`
	WalletAddress string    `json:"wallet_address"`
	CreatedAt     time.Time `json:"created_at"`
	ProcessedAt   time.Time `json:"processed_at"`
}
