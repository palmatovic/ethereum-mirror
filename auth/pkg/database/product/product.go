package product

import "time"

type Product struct {
	ProductId   int64     `json:"product_id" gorm:"column:ProductId;primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"column:Name"`
	Description string    `json:"description" gorm:"column:Description"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}
