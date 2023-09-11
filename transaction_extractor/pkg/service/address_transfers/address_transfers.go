package address_transfers

import (
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
	address_status_db "transaction-extractor/pkg/database/address_status"
	address_transfers_db "transaction-extractor/pkg/database/address_transfers"
	"transaction-extractor/pkg/model/address_transfers"
	"transaction-extractor/pkg/util"
)

func GetAddressTokenTransfers(db *gorm.DB, address string, address_status address_status_db.AddressStatus, browser playwright.Browser) (ats []address_transfers.AddressTransaction, err error) {

	var page playwright.Page
	defaultTimeout := float64(2000)
	page, err = browser.NewPage()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = page.Close()
	}()
	_, err = page.Goto("https://www.defined.fi/eth/" + address_status.TokenContractAddress)
	if err != nil {
		return nil, err
	}

	err = page.WaitForLoadState()
	if err != nil {
		return nil, err
	}

	page.WaitForTimeout(2000)
	//page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[3]/div/div/div[2]/div[1]/div/div")
	_ = page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[2]/div[5]").WaitFor(playwright.LocatorWaitForOptions{
		Timeout: &defaultTimeout,
	})
	expandTable := page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[2]/div[5]")
	err = expandTable.Click()
	if err != nil {
		return nil, err
	}
	page.WaitForTimeout(2000)

	_ = page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[3]/div/div/div[1]/div[7]/span/button").WaitFor(playwright.LocatorWaitForOptions{
		Timeout: &defaultTimeout,
	})
	filterButton := page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[3]/div/div/div[1]/div[7]/span/button")
	err = filterButton.Click()
	if err != nil {
		return nil, err
	}
	page.WaitForTimeout(2000)

	_ = page.Locator("xpath=//html/body/div[7]/div[3]/form/div[1]/div/input").WaitFor(playwright.LocatorWaitForOptions{
		Timeout: &defaultTimeout,
	})

	inputFiter := page.Locator("xpath=//html/body/div[7]/div[3]/form/div[1]/div/input")
	err = inputFiter.Fill(address)
	if err != nil {
		return nil, err
	}

	applyFilter := page.Locator("xpath=//html/body/div[7]/div[3]/form/div[2]/button[2]")
	err = applyFilter.Click()
	if err != nil {
		return nil, err
	}

	tableXpath := "xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[3]/div/div/div[2]/div[1]/div/div"
	page.WaitForTimeout(2000)

	_ = page.Locator(tableXpath).WaitFor(playwright.LocatorWaitForOptions{
		Timeout: &defaultTimeout,
	})

	table := page.Locator(tableXpath)
	rows, err := table.Locator("xpath=/div").All()
	if err != nil {
		return nil, err
	}

	for _, r := range rows {
		cols, err := r.Locator("xpath=/div/div").All()
		if err != nil {
			return nil, err
		}

		at := address_transfers.AddressTransaction{
			TxType:        "",
			Price:         "",
			Amount:        "",
			Total:         "",
			AgeTimestamp:  time.Time{},
			Asset:         address_status.TokenContractAddress,
			WalletAddress: address,
			CreatedAt:     time.Time{},
			ProcessedAt:   time.Time{},
		}
		for colNum, _ := range cols {

			if colNum == 0 {
				at.TxType, err = cols[colNum].InnerText()
				if err != nil {
					return nil, err
				}

			}

			if colNum == 1 {
				at.Price, err = cols[colNum].TextContent()
				if err != nil {
					return nil, err
				}
				at.Price = strings.ReplaceAll(at.Price, ",", "")
				at.Price = strings.ReplaceAll(at.Price, "$", "")

				util.CleanText(&at.Price)
			}

			if colNum == 2 {
				at.Amount, err = cols[colNum].TextContent()
				if err != nil {
					return nil, err
				}
				at.Amount = strings.ReplaceAll(at.Amount, ",", "")
				at.Amount = strings.ReplaceAll(at.Amount, "$", "")
				util.CleanText(&at.Amount)
			}

			if colNum == 3 {
				at.Total, err = cols[colNum].TextContent()
				if err != nil {
					return nil, err
				}
				at.Total = strings.ReplaceAll(at.Total, ",", "")
				at.Total = strings.ReplaceAll(at.Total, "$", "")
				util.CleanText(&at.Total)
			}

			if colNum == 4 {
				continue
			}

			if colNum == 5 {

				var colSpan = cols[colNum].Locator("xpath=/span")

				ageTimestamp, err := colSpan.GetAttribute("aria-label")
				if err != nil {
					return nil, err
				}
				print(ageTimestamp)

				at.AgeTimestamp, err = time.Parse("2006-01-02 15:04:05", ageTimestamp)
				if err != nil {
					return nil, err
				}
			}

		}
		ats = append(ats, at)
	}
	return ats, nil
}

func ScamCheck(tokenAddress string, browser playwright.Browser) (bool, error) {
	page, err := browser.NewPage()
	if err != nil {
		return false, err
	}

	_, err = page.Goto("https://gopluslabs.io/token-security/1/" + tokenAddress)
	if err != nil {
		return false, err
	}

	err = page.WaitForLoadState()
	if err != nil {
		return false, err
	}

	riskyItems := page.Locator("xpath=//html/body/div[1]/div[2]/div[2]/div[1]/div/div[3]/div[1]/div/div[2]")
	riskyItemsNum, err := riskyItems.TextContent()
	if err != nil {
		return false, err
	}

	attentionItems := page.Locator("xpath=//html/body/div[1]/div[2]/div[2]/div[1]/div/div[3]/div[1]/div/div[2]")
	attentionItemsNum, err := attentionItems.TextContent()
	if err != nil {
		return false, err
	}
	println(attentionItemsNum)

	err = page.Close()
	if err != nil {
		return false, err
	}
	if riskyItemsNum != "0" {
		return true, nil
	} else {
		return false, err
	}
}

func UpsertAddressTransfers(db *gorm.DB, addressTransfers []address_transfers.AddressTransaction) ([]address_transfers_db.Transaction, error) {
	var addressTransfersDb []address_transfers_db.Transaction
	for i := range addressTransfers {
		floatPrice, err := strconv.ParseFloat(addressTransfers[i].Price, 64)
		if err != nil {
			log.WithError(err).Errorf("Could not parse price %v with error: %v", addressTransfers[i], err)
		}
		floatAmount, err := strconv.ParseFloat(addressTransfers[i].Amount, 64)
		if err != nil {
			log.WithError(err).Errorf("Could not parse amount %v with error: %v", addressTransfers[i], err)
		}
		floatTotal, err := strconv.ParseFloat(addressTransfers[i].Total, 64)
		if err != nil {
			log.WithError(err).Errorf("Could not parse total %v with error: %v", addressTransfers[i], err)
		}
		addressTransfersDb = append(addressTransfersDb, address_transfers_db.Transaction{
			TxType:        addressTransfers[i].TxType,
			Price:         floatPrice,
			Amount:        floatAmount,
			Total:         floatTotal,
			AgeTimestamp:  addressTransfers[i].AgeTimestamp,
			Asset:         addressTransfers[i].Asset,
			WalletAddress: addressTransfers[i].WalletAddress,
		})
	}
	for j := range addressTransfersDb {
		//var asd address_transfers_db.Transaction
		var err error
		//grosso problema, mi devo tenere un univoco della transazione e qua non ho il transaction hash...
		// risolto, entro su action e me prendo il thHash
		if err = db.Create(&addressTransfersDb[j]).Error; err != nil {
			return nil, err
		}

		//err = db.Where("AddressId = ? AND TokenContractAddress = ?", address, addressStatusesDb[j].TokenContractAddress).First(&asd).Error
		//if err != nil {
		//	if errors.Is(gorm.ErrRecordNotFound, err) {
		//		if err = db.Create(&addressStatusesDb[j]).Error; err != nil {
		//			return nil, err
		//		}
		//		// inserire nella la tabella delle transazioni che hanno portato a quel token amount x il token contract address e l'addressId (etherscan puoi filtrare per contratto e holder wallet) (che non esiste)
		//
		//	} else {
		//		return nil, err
		//	}
		//} else {
		//	if asd.TokenAmount != addressStatusesDb[j].TokenAmount {
		//		// aggiorna la tabelle address_status con il nuovo token amont e aggiorna la tabella delle transazioni con solo le transazioni nuove (che non esiste vedi sopra)
		//	}
		//}
	}
	return nil, nil
}
