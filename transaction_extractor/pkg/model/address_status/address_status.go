package address_status

type AddressStatus struct {
	Contract string  `json:"contract"`
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Amount   float64 `json:"amount"`
}
