package sync

import (
	"errors"
	"github.com/playwright-community/playwright-go"
	"gorm.io/gorm"
	address_db "transaction-extractor/pkg/database/address"
	address_status_db "transaction-extractor/pkg/database/address_status"
	address_status_model "transaction-extractor/pkg/model/address_status"
	address_status_service "transaction-extractor/pkg/service/address_status"
)

type Env struct {
	playwright.Browser
	Database  *gorm.DB
	Addresses []string
}

func (e *Env) SyncTransactions() (response interface{}, err error) {
	for _, address := range e.Addresses {
		var addressRecord address_db.Address
		if err = e.Database.Where("AddressId = ?", address).First(&addressRecord).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				addressRecord = address_db.Address{AddressId: address}
				if errCreate := e.Database.Create(&addressRecord).Error; errCreate != nil {
					return nil, errCreate
				}
			} else {
				return nil, err
			}
		}
		var addressStatuses []address_status_model.AddressStatus
		addressStatuses, err = address_status_service.GetAddressStatus(e.Browser, address)
		if err != nil {
			return nil, err
		}
		if len(addressStatuses) > 0 {
			var savedAddressStatuses []address_status_db.AddressStatus
			savedAddressStatuses, err = address_status_service.SaveNewAddressStatus(e.Database, address, addressStatuses)
			if err != nil {
				return nil, err
			}
			if len(savedAddressStatuses) > 0 {

			}
			// get tokens transfer for wallet and save to transactions table
		}
	}
	return nil, nil
}
