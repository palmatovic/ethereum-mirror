package model

import "time"

type Transaction struct {
	TransactionHash      string
	TxType               string
	Method               string
	Block                string
	AgeMillis            string
	AgeTimestamp         time.Time
	AgeDistanceFromQuery string
	GasPrice             string
	From                 string
	To                   string
	Amount               string
	Value                string
	Asset                string
	TxnFee               string
}
