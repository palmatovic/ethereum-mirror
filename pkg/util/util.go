package util

import (
	"github.com/playwright-community/playwright-go"
	"strings"
)

type Selector string

const (
	TxHash   Selector = "#referralLink-1"
	TxAction Selector = "xpath=/html/body/main/div/div[1]/div[1]/div/div[1]/div[5]/div[2]/div/div/div"
	TxStatus Selector = "xpath=/html/body/main/div/div[1]/div[1]/div/div[1]/div[2]/div[2]/span"
)

func CleanText(text *string) {
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
