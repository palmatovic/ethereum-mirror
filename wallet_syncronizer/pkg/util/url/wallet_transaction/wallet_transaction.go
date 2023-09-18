package wallet_transaction

import (
	"wallet-syncronizer/pkg/util/url"
)

const Resource = url.BaseUrl + "/wallet-transaction"

const (
	Get     = Resource
	GetList = Resource + "/list"
)
