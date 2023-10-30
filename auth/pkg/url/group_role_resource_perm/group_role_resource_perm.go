package group_role_resource_perm

import "auth/pkg/url"

const Resource = url.BaseUrl + "/group-role-resource-perm"

const (
	GraphQL = Resource + "/graphql"
	Get     = Resource + "/get/:" + Id
	List    = Resource + "/list"
	Create  = Resource
	Update  = Resource
	Delete  = Resource + "/remove/:" + Id
)

const Id = "group_role_resource_perm_id"
