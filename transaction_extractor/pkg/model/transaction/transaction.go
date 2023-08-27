package transaction

import "time"

type Transaction struct {
	TxHash               string    `json:"tx_hash"`
	TxType               string    `json:"tx_type"`
	Method               string    `json:"method"`
	Block                string    `json:"block"`
	AgeTimestamp         time.Time `json:"age_timestamp"`
	AgeDistanceFromQuery string    `json:"age_distance_from_query"`
	AgeMillis            string    `json:"age_millis"`
	From                 string    `json:"from"`
	To                   string    `json:"to"`
	Amount               string    `json:"amount"`
	Value                string    `json:"value"`
	Asset                string    `json:"asset"`
	TxnFee               string    `json:"txn_fee"`
	GasPrice             string    `json:"gas_price"`
	WalletAddress        string    `json:"wallet_address"`
}
