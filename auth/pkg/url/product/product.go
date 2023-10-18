package product

import "auth/pkg/url"

const Resource = url.BaseUrl + "/product"

const (
	GraphQL = Resource + "/graphql"
	Get     = Resource + "/get/:" + string(Id)
	List    = Resource + "/list"
	Create  = Resource
	Update  = Resource
	Delete  = Resource + "/remove/:" + string(Id)
)

type Parameter string

const Id Parameter = "product_id"
