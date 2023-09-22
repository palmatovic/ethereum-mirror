package wallet_transaction

import (
	"wallet-syncronizer/pkg/util/url"
)

const Resource = url.BaseUrl + "/wallet-transaction"

const (
	Get  = Resource + "/:" + string(Id)
	List = Resource + "/list"
)

type Parameter string

const Id Parameter = "wallet_transaction_id"
