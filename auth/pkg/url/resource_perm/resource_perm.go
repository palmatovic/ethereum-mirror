package resource_perm

import "auth/pkg/url"

const Resource = url.BaseUrl + "/resource-perm"

const (
	GraphQL = Resource + "/graphql"
	Get     = Resource + "/get/:" + Id
	List    = Resource + "/list"
	Create  = Resource
	Update  = Resource
	Delete  = Resource + "/remove/:" + Id
)

const Id = "resource_perm_id"
