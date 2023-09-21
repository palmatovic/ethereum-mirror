package wallet

import (
	"wallet-syncronizer/pkg/util/url"
)

const Resource = url.BaseUrl + "/wallet"

const (
	Get    = Resource + "/:" + string(Id)
	List   = Resource + "/list"
	Create = Resource
	Update = Resource
	Delete = Resource + "/:" + string(Id)
)

type Parameter string

const Id Parameter = "wallet_id"
