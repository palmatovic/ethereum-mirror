package wallet_token

import (
	"wallet-synchronizer/pkg/util/url"
	token_url "wallet-synchronizer/pkg/util/url/token"
	wallet_url "wallet-synchronizer/pkg/util/url/wallet"
)

const Resource = url.BaseUrl + "/wallet-token"

const (
	Get  = Resource + "/:" + string(wallet_url.Id) + "/:" + string(token_url.Id)
	List = Resource + "/list"
)
