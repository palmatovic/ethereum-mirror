package user_group_role

import (
	"auth/pkg/database/group_role"
	"auth/pkg/database/user"
	"time"
)

type UserGroupRole struct {
	UserGroupRoleId int64 `json:"user_group_role_id" gorm:"column:UserGroupRoleId;primaryKey;autoIncrement"`
	UserId          int64 `json:"user_id" gorm:"column:UserId;uniqueIndex:UserIdGroupRoleIdIdx"`
	user.User
	GroupRoleId int64 `json:"group_role_id" gorm:"column:GroupRoleId;uniqueIndex:UserIdGroupRoleIdIdx"`
	group_role.GroupRole
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (UserGroupRole) TableName() string {
	return "UserGroupRole"
}
