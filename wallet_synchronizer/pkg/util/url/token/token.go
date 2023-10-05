package token

import (
	"wallet-synchronizer/pkg/util/url"
)

const Resource = url.BaseUrl + "/token"

const (
	Get  = Resource + "/:" + string(Id)
	List = Resource + "/list"
)

type Parameter string

const Id Parameter = "token_id"
