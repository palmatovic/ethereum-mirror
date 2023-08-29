package address_status

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"strings"
	"transaction-extractor/pkg/model/address_status"
	"transaction-extractor/pkg/util"
)

// GetAddressStatus returns a list of all token balance by address_status
func GetAddressStatus(browser playwright.Browser, address string) (addressStatuses []address_status.AddressStatus, err error) {

	page, err := browser.NewPage()
	if err != nil {
		return nil, err
	}
	page.SetDefaultTimeout(1000 * 40)
	defer func() {
		_ = page.Close()
	}()

	_, err = page.Goto(fmt.Sprintf("https://etherscan.io/tokenholdings?a=%s", address))
	if err != nil {
		return nil, err
	}

	//_, err = page.WaitForSelector(string(util.CloudFlare))
	//if err != nil {
	//	return nil, err
	//}

	//
	//cloudFlare, err := page.QuerySelector(string(util.CloudFlare))
	//if err != nil {
	//	return nil, err
	//}
	//
	//if err = cloudFlare.Check(); err != nil {
	//	return nil, err
	//}

	_, err = page.WaitForSelector(string(util.TableBody))
	if err != nil {
		return nil, err
	}

	tableBody, err := page.QuerySelector(string(util.TableBody))
	if err != nil {
		return nil, err
	}

	rows, err := tableBody.QuerySelectorAll(string(util.RelativeTableRows))
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		cells, err := row.QuerySelectorAll(string(util.RelativeTableData))
		if err != nil {
			return nil, err
		}

		var cellData []string
		for i := range cells {
			text, err := cells[i].TextContent()
			if err != nil {
				return nil, err
			}
			text = strings.Join(strings.Fields(text), " ")
			util.CleanText(&text)
			if len(text) == 0 {
				continue
			}
			//if slices.Contains([]int{7, 9}, i) {
			//	el, err := cells[i].QuerySelector(`xpath=div/a[@class="js-clipboard link-secondary "]`)
			//	if err != nil {
			//		return nil, err
			//	}
			//	text, err = el.GetAttribute("data-clipboard-text")
			//	if err != nil {
			//		return nil, err
			//	}
			//	util.CleanText(&text)
			//}
			cellData = append(cellData, text)
			println(cellData)
		}

	}
	return addressStatuses, nil
}

//
//// SaveNewAddressStatus saves only new address_status data in the database AddressStatus table
//func SaveNewAddressStatus(database *gorm.DB, transactions []address_status.AddressStatus) ([]address_status_db.AddressStatus, error) {
//
//}
