package group_role

import "auth/pkg/url"

const Resource = url.BaseUrl + "/group-role"

const (
	GraphQL = Resource + "/graphql"
	Get     = Resource + "/get/:" + Id
	List    = Resource + "/list"
	Create  = Resource
	Update  = Resource
	Delete  = Resource + "/remove/:" + Id
)

const Id = "group_role_id"
