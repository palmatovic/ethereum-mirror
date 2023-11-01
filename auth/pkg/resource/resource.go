package resource

const (
	Login                 = "login"
	Logout                = "logout"
	Otp                   = "otp"
	Product               = "product"
	Company               = "company"
	Group                 = "group"
	Role                  = "role"
	GroupRole             = "group_role"
	Resource              = "resource"
	Perm                  = "perm"
	ResourcePerm          = "resource_perm"
	GroupRoleResourcePerm = "group_role_resource_perm"
	User                  = "user"
	UserGroupRole         = "user_group_role"
	UserProduct           = "user_product"
	UserResourcePerm      = "user_resource_perm"
)

var AllResource = []string{
	Login,
	Logout,
	Otp,
	Product,
	Company,
	Group,
	Role,
	GroupRole,
	Resource,
	Perm,
	ResourcePerm,
	GroupRoleResourcePerm,
	User,
	UserGroupRole,
	UserProduct,
	UserResourcePerm,
}
