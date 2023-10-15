package user_resource_perm

import (
	"auth/pkg/database/resource_perm"
	"auth/pkg/database/user"
	"time"
)

type UserResourcePerm struct {
	UserResourcePermId int `json:"user_resource_id" gorm:"UserResourcePermId;autoIncrement;primaryKey"`
	UserId             int `json:"user_id" gorm:"column:UserId;uniqueIndex:UserResourcePermIdx"`
	user.User
	ResourcePermId int `json:"resource_perm_id" gorm:"ResourcePermId;uniqueIndex:UserResourcePermIdx"`
	resource_perm.ResourcePerm
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (UserResourcePerm) TableName() string {
	return "UserResourcePerm"
}
