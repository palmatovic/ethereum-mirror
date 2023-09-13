package goplus

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func ScamCheck(tokenAddress string, browser playwright.Browser) (risk int, warning int, err error) {
	var page playwright.Page
	page, err = browser.NewPage()
	if err != nil {
		return 0, 0, err
	}

	defer func() {
		_ = page.Close()
	}()

	_, err = page.Goto(fmt.Sprintf("https://gopluslabs.io/token-security/1/%s", tokenAddress))
	if err != nil {
		return 0, 0, err
	}

	err = page.WaitForLoadState()
	if err != nil {
		return 0, 0, err
	}

	riskItemsLocator := page.Locator("xpath=//html/body/div[1]/div[2]/div[2]/div[1]/div/div[3]/div[1]/div/div[2]")
	riskItemStr, err := riskItemsLocator.TextContent()
	if err != nil {
		return 0, 0, err
	}
	riskItems, err := strconv.Atoi(riskItemStr)
	if err != nil {
		return 0, 0, err
	}

	attentionItemsLocator := page.Locator("xpath=//html/body/div[1]/div[2]/div[2]/div[1]/div/div[3]/div[1]/div/div[2]")
	attentionItemsStr, err := attentionItemsLocator.TextContent()
	if err != nil {
		return 0, 0, err
	}
	attentionItems, err := strconv.Atoi(attentionItemsStr)
	if err != nil {
		return 0, 0, err
	}
	logrus.WithFields(logrus.Fields{"contract_address": tokenAddress, "risk": riskItems, "waning": attentionItems}).Info("scam check data retrieved")
	return riskItems, attentionItems, nil

}
