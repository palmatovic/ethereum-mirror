package company

import "time"

type Company struct {
	CompanyId int64     `json:"company_id" gorm:"column:CompanyId;primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"column:Name"`
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (Company) TableName() string {
	return "Company"
}
