package token

import (
	"wallet-syncronizer/pkg/util/url"
)

const Resource = url.BaseUrl + "/token"

const (
	Get     = Resource + "/:" + string(Id)
	GetList = Resource + "/list"
)

type Parameter string

const Id Parameter = "token_id"
