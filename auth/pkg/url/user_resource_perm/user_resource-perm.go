package user_resource_perm

import "auth/pkg/url"

const Resource = url.BaseUrl + "/user-resource-perm"

const (
	GraphQL = Resource + "/graphql"
	Get     = Resource + "/get/:" + string(Id)
	List    = Resource + "/list"
	Create  = Resource
	Update  = Resource
	Delete  = Resource + "/remove/:" + string(Id)
)

type Parameter string

const Id Parameter = "user_resource_perm_id"
