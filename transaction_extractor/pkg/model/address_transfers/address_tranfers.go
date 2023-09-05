package address_transfers

import "time"

type AddressTransaction struct {
	TxType        string    `json:"tx_type"`
	Price         string    `json:"price"`
	Amount        string    `json:"amount"`
	Total         string    `json:"total"`
	AgeTimestamp  time.Time `json:"age_timestamp"`
	Asset         string    `json:"asset"`
	WalletAddress string    `json:"wallet_address"`
	CreatedAt     time.Time `json:"created_at"`
	ProcessedAt   time.Time `json:"processed_at"`
}
