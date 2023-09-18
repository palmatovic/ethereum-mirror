package token

import (
	"wallet-syncronizer/pkg/util/url"
)

const Resource = url.BaseUrl + "/token"

const (
	Get     = Resource
	GetList = Resource + "/list"
)
