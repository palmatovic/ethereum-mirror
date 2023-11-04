package util

const (
	WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_EXPIRED        = "expired"
	WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_OPENED         = "open"           // la transazione buy, inserita diventa opened e viene monitorata
	WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_CLOSED         = "closed"         // la transazione è stata processata correttamente
	WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_REGISTERED     = "registered"     // appena segnalata, la transazione viene creata come registered
	WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_ALREADY_CLOSED = "already_closed" // transazione buy che viene chiusa perchè non ci sono più token disponibili (magari a causa di una transazione sell eseguita prima)
)

const (
	WALLET_TRANSACTION_TYPE_SELL = "Sell"
	WALLET_TRANSACTION_TYPE_BUY  = "Buy"
)
