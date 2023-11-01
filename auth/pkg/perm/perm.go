package perm

const (
	Create = "create"
	Update = "update"
	Get    = "get"
	Delete = "delete"
	List   = "list"
)

var AllPerm = []string{
	Create,
	Update,
	Get,
	Delete,
	List,
}
