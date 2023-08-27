package util

import (
	"github.com/playwright-community/playwright-go"
	"strings"
)

type Selector string

const (
	TxHash Selector = "#referralLink-1"
	//TxAction                  Selector = "xpath=/html/body/main/div/div[1]/div[1]/div/div[1]/div[5]/div[2]/div/div/div"
	TxAction                  Selector = "#wrapperContent"
	TxStatus                  Selector = "xpath=/html/body/main/div/div[1]/div[1]/div/div[1]/div[2]/div[2]/span"
	TxBlock                   Selector = "xpath=/html/body/main/div/div[1]/div[1]/div/div[1]/div[3]/div[2]/div/span[1]/a"
	TxTimestamp               Selector = "xpath=/html/body/main/div/div[1]/div[1]/div/div[1]/div[4]/div/div[2]"
	TxFrom                    Selector = "xpath=/html/body/main/div/div[1]/div[1]/div/div[1]/div[7]/div[2]/div/span/a"
	TxInteractedWithToSuccess Selector = "#ContentPlaceHolder1_maintable > div.card.p-5.mb-3 > div:nth-child(11) > div.col-md-9 > div > span"
	TxInteractedWithToFail    Selector = "#ContentPlaceHolder1_maintable > div.card.p-5.mb-3 > div:nth-child(9) > div.col-md-9 > div:nth-child(1) > span.me-1"
)

const (
	HeaderPageSize  Selector = "xpath=/html/body/main/form/main/section[2]/div[3]/div[3]/div/select"
	HeaderTable     Selector = "table.table tbody tr"
	HeaderTableData Selector = "td"
)

func CleanText(text *string) {
	*text = strings.TrimSpace(*text)
	*text = strings.TrimPrefix(*text, " ")
	*text = strings.TrimPrefix(*text, "\n")
	*text = strings.TrimPrefix(*text, "\t")
}

func GetObjectByPage(page playwright.Page, dSelector Selector) (string, error) {
	var err error
	selector := string(dSelector)
	_, err = page.WaitForSelector(selector)
	if err != nil {
		return "", err
	}

	object, err := page.QuerySelector(selector)
	if err != nil {
		return "", err
	}

	text, err := object.TextContent()
	if err != nil {
		return "", err
	}

	CleanText(&text)

	return text, nil
}
