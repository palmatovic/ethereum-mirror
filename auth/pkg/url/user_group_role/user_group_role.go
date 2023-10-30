package user_group_role

import "auth/pkg/url"

const Resource = url.BaseUrl + "/user-group-role"

const (
	GraphQL = Resource + "/graphql"
	Get     = Resource + "/get/:" + Id
	List    = Resource + "/list"
	Create  = Resource
	Update  = Resource
	Delete  = Resource + "/remove/:" + Id
)

const Id = "user_group_role_id"
