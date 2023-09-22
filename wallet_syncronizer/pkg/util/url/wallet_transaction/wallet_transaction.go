package wallet_transaction

import (
	"wallet-syncronizer/pkg/util/url"
	token_url "wallet-syncronizer/pkg/util/url/token"
	wallet_token_url "wallet-syncronizer/pkg/util/url/wallet_token"
)

const Resource = url.BaseUrl + "/wallet-transaction"

const (
	Get     = Resource + "/:" + string(wallet_token_url.Id) + "/:" + string(token_url.Id)
	GetList = Resource + "/list"
)
