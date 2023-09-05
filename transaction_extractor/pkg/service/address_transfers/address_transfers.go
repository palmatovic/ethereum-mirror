package address_transfers

import (
	"github.com/playwright-community/playwright-go"
	"gorm.io/gorm"
	"strings"
	"time"
	address_status_db "transaction-extractor/pkg/database/address_status"
	"transaction-extractor/pkg/model/address_transfers"
)

func GetAddressTokenTransfers(db *gorm.DB, address string, address_status address_status_db.AddressStatus, browser playwright.Browser) {
	page, err := browser.NewPage()
	if err != nil {
		return
	}

	_, err = page.Goto("https://www.defined.fi/eth/" + address_status.TokenContractAddress)
	if err != nil {
		return
	}

	err = page.WaitForLoadState()
	if err != nil {
		return
	}
	page.WaitForTimeout(2000)

	//page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[3]/div/div/div[2]/div[1]/div/div")

	expandTable := page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[2]/div[5]")
	err = expandTable.Click()
	if err != nil {
		return
	}

	page.WaitForTimeout(2000)

	filterButton := page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[3]/div/div/div[1]/div[7]/span/button")
	err = filterButton.Click()
	if err != nil {
		return
	}
	page.WaitForTimeout(2000)

	inputFiter := page.Locator("xpath=//html/body/div[7]/div[3]/form/div[1]/div/input")
	err = inputFiter.Fill(address)
	if err != nil {
		return
	}

	page.WaitForTimeout(2000)

	applyFilter := page.Locator("xpath=//html/body/div[7]/div[3]/form/div[2]/button[2]")
	err = applyFilter.Click()
	if err != nil {
		return
	}

	tableXpath := "xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[3]/div/div/div[2]/div[1]/div/div"

	table := page.Locator(tableXpath)
	rows, err := table.Locator("xpath=/div").All()
	if err != nil {
		return
	}

	for _, r := range rows {
		cols, err := r.Locator("xpath=/div/div").All()
		if err != nil {
			return
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
					return
				}
			}

			if colNum == 1 {
				at.Price, err = cols[colNum].InnerText()
				if err != nil {
					return
				}
				at.Price = strings.TrimSpace(at.Price[1:])
			}

			if colNum == 2 {
				at.Amount, err = cols[colNum].InnerText()
				if err != nil {
					return
				}
			}

			if colNum == 3 {
				at.Total, err = cols[colNum].InnerText()
				if err != nil {
					return
				}
			}

			if colNum == 4 {
				continue
			}

			if colNum == 5 {
				var colSpan = cols[colNum].Locator("xpath=/span")

				ageTimestamp, err := colSpan.GetAttribute("aria-label")
				if err != nil {
					return
				}
				print(ageTimestamp)

				at.AgeTimestamp, err = time.Parse("2006-01-02 15:04:05", ageTimestamp)
				if err != nil {
					return
				}
			}

		}
	}

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

	err := page.Close()
	if err != nil {
		return false, err
	}
	if riskyItemsNum != "0" {
		return true, nil
	} else {
		return false, err
	}
}
