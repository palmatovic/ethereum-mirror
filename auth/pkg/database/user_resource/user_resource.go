package user_resource

import (
	"auth/pkg/database/resource"
	"auth/pkg/database/user"
	"time"
)

type UserResource struct {
	UserResourceId int `json:"user_resource_id" gorm:"UserResourceId;autoIncrement;primaryKey"`
	UserId         int `json:"user_id" gorm:"column:UserId;uniqueIndex:UserResourceIdx"`
	user.User
	ResourceId string `json:"resource_id" gorm:"ResourceId;uniqueIndex:UserResourceIdx"`
	resource.Resource
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (UserResource) TableName() string {
	return "UserResource"
}
