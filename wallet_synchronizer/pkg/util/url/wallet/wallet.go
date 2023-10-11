package wallet

import (
	"wallet-synchronizer/pkg/util/url"
)

const Resource = url.BaseUrl + "/wallet"

const (
	GraphQL = Resource + "/graphql"
	Get     = Resource + "/get/:" + string(Id)
	List    = Resource + "/list"
	Create  = Resource
	Update  = Resource
	Delete  = Resource + "/remove/:" + string(Id)
)

type Parameter string

const Id Parameter = "wallet_id"
