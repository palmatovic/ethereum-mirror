package wallet_token

import (
	"wallet-syncronizer/pkg/util/url"
)

const Resource = url.BaseUrl + "/wallet-token"

const (
	Get  = Resource + "/:" + string(Id)
	List = Resource + "/list"
)

type Parameter string

const Id Parameter = ":wallet_token_id"
