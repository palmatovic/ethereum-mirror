package group_role

import (
	"auth/pkg/database/group"
	"auth/pkg/database/role"
	"time"
)

type GroupRole struct {
	GroupRoleId int64 `json:"group_role_id" gorm:"column:GroupRoleId;primaryKey;autoIncrement"`
	GroupId     int64 `json:"group_id" gorm:"column:GroupId;uniqueIndex:GroupIdRoleIdIdx"`
	group.Group
	RoleId int64 `json:"role_id" gorm:"column:RoleId;uniqueIndex:GroupIdRoleIdIdx"`
	role.Role
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (GroupRole) TableName() string {
	return "GroupRole"
}
