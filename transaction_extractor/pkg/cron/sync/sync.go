package sync

import (
	"errors"
	"github.com/playwright-community/playwright-go"
	"gorm.io/gorm"
	address_db "transaction-extractor/pkg/database/address"
	address_status_db "transaction-extractor/pkg/database/address_status"
	address_status_model "transaction-extractor/pkg/model/address_status"
	address_status_service "transaction-extractor/pkg/service/address_status"
	address_transfers_service "transaction-extractor/pkg/service/address_transfers"
)

type Env struct {
	playwright.Browser
	Database      *gorm.DB
	Addresses     []string
	AlchemyApiKey string
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
		addressStatuses, err = address_status_service.GetAddressStatus(e.Database, e.Browser, address, e.AlchemyApiKey)
		if err != nil {
			return nil, err
		}
		if len(addressStatuses) > 0 {
			var savedAddressStatuses []address_status_db.AddressStatus
			savedAddressStatuses, err = address_status_service.UpsertAddressStatus(e.Database, address, addressStatuses)
			if err != nil {
				return nil, err
			}
			if len(savedAddressStatuses) > 0 {

				//GET TOKEN TRANSFER FOR ADDRESS

				for _, sas := range savedAddressStatuses {
					//address_transfers_service.GetAddressTokenTransfers(e.Database, address, sas, e.Browser)
					check, err := address_transfers_service.ScamCheck(sas.TokenContractAddress, e.Browser)
					if err != nil {
						return nil, err
					}
					println(sas.TokenContractAddress)
					println(check)
				}
			}

			// get tokens transfer for wallet and save to transactions table
		}
	}
	return nil, nil
}
