package perm

import "auth/pkg/url"

const Resource = url.BaseUrl + "/perm"

const (
	GraphQL = Resource + "/graphql"
	Get     = Resource + "/get/:" + Id
	List    = Resource + "/list"
	Create  = Resource
	Update  = Resource
	Delete  = Resource + "/remove/:" + Id
)

const Id = "perm_id"
