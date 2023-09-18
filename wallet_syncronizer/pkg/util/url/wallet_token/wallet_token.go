package wallet_token

import (
	"wallet-syncronizer/pkg/util/url"
)

const Resource = url.BaseUrl + "/wallet-token"

const (
	Get     = Resource
	GetList = Resource + "/list"
)
