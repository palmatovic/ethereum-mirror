package role

import "auth/pkg/url"

const Resource = url.BaseUrl + "/role"

const (
	GraphQL = Resource + "/graphql"
	Get     = Resource + "/get/:" + Id
	List    = Resource + "/list"
	Create  = Resource
	Update  = Resource
	Delete  = Resource + "/remove/:" + Id
)

const Id = "role_id"
