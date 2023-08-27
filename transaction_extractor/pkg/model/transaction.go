package model

import "time"

type Transaction struct {
	TransactionHash      string    `json:"transaction_hash"`
	Method               string    `json:"method"`
	Block                string    `json:"block"`
	AgeMillis            string    `json:"age_millis"`
	AgeTimestamp         time.Time `json:"age_timestamp"`
	AgeDistanceFromQuery string    `json:"age_distance_from_query"`
	GasPrice             string    `json:"gas_price"`
	From                 string    `json:"from"`
	To                   string    `json:"to"`
	InOut                string    `json:"in_out"`
	Value                string    `json:"value"`
	TxnFee               string    `json:"txn_fee"`
}
