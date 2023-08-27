package transaction_detail

type TransactionDetail struct {
	TransactionHash      string `json:"transaction_hash"`
	TransactionStatus    string `json:"transaction_status"`
	TransactionBlock     string `json:"transaction_block"`
	TransactionTimestamp string `json:"transaction_timestamp"`
	TransactionAction    string `json:"transaction_action"`
	TransactionFrom      string `json:"transaction_from"`
	TransactionTo        string `json:"transaction_to"`
}
