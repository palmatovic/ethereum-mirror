package product

import "time"

type Product struct {
	ProductId                         int64     `json:"product_id" gorm:"column:ProductId;primaryKey;autoIncrement"`
	Name                              string    `json:"name" gorm:"column:Name;uniqueIndex:NameIdx"`
	Description                       string    `json:"description" gorm:"column:Description"`
	ServerCert                        []byte    `json:"server_cert" gorm:"column:ServerCert"`
	ServerKey                         []byte    `json:"server_key" gorm:"column:ServerKey"`
	CaCert                            []byte    `json:"ca_cert" gorm:"column:CaCert"`
	CaKey                             []byte    `json:"ca_key" gorm:"column:CaKey"`
	SSLExpired                        bool      `json:"ssl_expired" gorm:"column:SSLExpired;default=0"`
	RSAPrivateKey                     []byte    `json:"rsa_private_key" gorm:"column:RSAPrivateKey"`
	RSAPublicKey                      []byte    `json:"rsa_public_key" gorm:"column:RSAPublicKey"`
	AccessTokenExpiresInMinutes       int64     `json:"access_token_expires_in_minutes" gorm:"column:AccessTokenExpiresInMinutes"`
	RefreshTokenExpiresInMinutes      int64     `json:"refresh_token_expires_in_minutes" gorm:"column:RefreshTokenExpiresInMinutes"`
	RSAExpirationDate                 time.Time `json:"jwt_expiration_date" gorm:"column:JwtExpirationDate"`
	RSAExpired                        bool      `json:"rsa_expired" gorm:"column:RSAExpired;default=0"`
	AES256EncryptionKey               []byte    `json:"aes_256_encryption_key" gorm:"column:AES256EncryptionKey"`
	AES256EncryptionKeyExpirationDate time.Time `json:"aes_256_encryption_key_expiration_date" gorm:"column:AES256EncryptionKeyExpirationDate"`
	AES256Expired                     bool      `json:"aes_256_expired" gorm:"column:AES245Expired;default=0"`
	CreatedAt                         time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt                         time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (Product) TableName() string {
	return "Product"
}
