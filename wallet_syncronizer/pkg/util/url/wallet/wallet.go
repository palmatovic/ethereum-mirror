package wallet

import (
	"wallet-syncronizer/pkg/util/url"
)

const Resource = url.BaseUrl + "/wallet"

const (
	Get     = Resource + "/:wallet_id"
	GetList = Resource + "/list"
	Create  = Resource
	Update  = Resource
	Delete  = Resource + "/:wallet_id"
)
