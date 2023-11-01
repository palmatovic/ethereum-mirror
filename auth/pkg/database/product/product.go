package product

import "time"

type Product struct {
	ProductId   int64  `json:"product_id" gorm:"column:ProductId;primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"column:Name;uniqueIndex:ProductNameIdx"`
	Description string `json:"description" gorm:"column:Description"`

	ServerCert        string    `json:"server_cert" gorm:"column:ServerCert"`
	ServerKey         string    `json:"server_key" gorm:"column:ServerKey"`
	CaCert            string    `json:"ca_cert" gorm:"column:CaCert"`
	CaKey             string    `json:"ca_key" gorm:"column:CaKey"`
	SSLExpirationDate time.Time `json:"ssl_expiration_date" gorm:"column:SSLExpirationDate"`
	SSLExpired        bool      `json:"ssl_expired" gorm:"column:SSLExpired;default=0"`

	RSAPrivateKey                string    `json:"rsa_private_key" gorm:"column:RSAPrivateKey"`
	RSAPublicKey                 string    `json:"rsa_public_key" gorm:"column:RSAPublicKey"`
	AccessTokenExpiresInMinutes  int64     `json:"access_token_expires_in_minutes" gorm:"column:AccessTokenExpiresInMinutes"`
	RefreshTokenExpiresInMinutes int64     `json:"refresh_token_expires_in_minutes" gorm:"column:RefreshTokenExpiresInMinutes"`
	RSAExpirationDate            time.Time `json:"jwt_expiration_date" gorm:"column:RSAExpirationDate"`
	RSAExpired                   bool      `json:"rsa_expired" gorm:"column:RSAExpired;default=0"`

	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (Product) TableName() string {
	return "Product"
}
