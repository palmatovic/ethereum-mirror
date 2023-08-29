package address

func (Address) TableName() string {
	return "Address"
}

type Address struct {
	AddressId string `gorm:"column:AddressId;primaryKey;varchar(1024)"`
}
