package token

import (
	"wallet-synchronizer/pkg/url"
)

const Resource = url.BaseUrl + "/token"

const (
	GraphQL = Resource + "/graphql"
	Get     = Resource + "/get/:" + string(Id)
	List    = Resource + "/list"
)

type Parameter string

const Id Parameter = "token_id"
