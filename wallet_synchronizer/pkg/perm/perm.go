package perm

const (
	GraphQL = "graphql"
	Create  = "create"
	Update  = "update"
	Get     = "get"
	Delete  = "delete"
	List    = "list"
)

var AllPerm = []string{
	GraphQL,
	Create,
	Update,
	Get,
	Delete,
	List,
}
