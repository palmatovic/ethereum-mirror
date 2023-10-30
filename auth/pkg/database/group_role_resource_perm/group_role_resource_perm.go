package group_role_resource_perm

import (
	"auth/pkg/database/group_role"
	"auth/pkg/database/resource_perm"
	"time"
)

type GroupRoleResourcePerm struct {
	GroupRoleResourcePermId int64                      `json:"group_role_resource_perm_id" gorm:"column:GroupRoleResourcePermId;primaryKey;autoIncrement"`
	GroupRoleId             int64                      `json:"group_role_id" gorm:"column:GroupRoleId;uniqueIndex:GroupRoleIdResourcePermIdIdx"`
	GroupRole               group_role.GroupRole       `json:"-"`
	ResourcePermId          int64                      `json:"resource_perm_id" gorm:"column:ResourcePermId;uniqueIndex:GroupRoleIdResourcePermIdIdx"`
	ResourcePerm            resource_perm.ResourcePerm `json:"-"`
	CreatedAt               time.Time                  `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt               time.Time                  `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (GroupRoleResourcePerm) TableName() string {
	return "GroupRoleResourcePerm"
}
