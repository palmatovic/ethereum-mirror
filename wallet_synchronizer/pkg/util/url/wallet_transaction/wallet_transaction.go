package wallet_transaction

import (
	"wallet-synchronizer/pkg/util/url"
)

const Resource = url.BaseUrl + "/wallet-transaction"

const (
	GraphQL = Resource + "/graphql"
	Get     = Resource + "/get/:" + string(Id)
	List    = Resource + "/list"
)

type Parameter string

const Id Parameter = "wallet_transaction_id"
