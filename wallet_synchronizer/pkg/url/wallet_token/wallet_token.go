package wallet_token

import (
	"wallet-synchronizer/pkg/url"
	token_url "wallet-synchronizer/pkg/url/token"
	wallet_url "wallet-synchronizer/pkg/url/wallet"
)

const Resource = url.BaseUrl + "/wallet-token"

const (
	GraphQL = Resource + "/graphql"
	Get     = Resource + "/get/:" + string(wallet_url.Id) + "/:" + string(token_url.Id)
	List    = Resource + "/list"
)
