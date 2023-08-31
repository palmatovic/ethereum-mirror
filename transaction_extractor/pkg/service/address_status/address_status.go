package address_status

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"sync"
	"transaction-extractor/pkg/model/address_status"
	"transaction-extractor/pkg/util"
)

type TokenBalance struct {
	ContractAddress string `json:"contractAddress"`
	TokenBalance    string `json:"tokenBalance"`
}

type Result struct {
	Address       string         `json:"address"`
	TokenBalances []TokenBalance `json:"tokenBalances"`
}

type GetTokenBalanceResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  Result `json:"result"`
}

// GetAddressStatus returns a list of all token balance by address_status
func GetAddressStatus(browser playwright.Browser, address string) (addressStatuses []address_status.AddressStatus, err error) {

	// Alchemy URL
	const baseURL = "https://eth-mainnet.g.alchemy.com/v2/owUCVigVvnHA63o0C6mh3yrf3jxMkV7b"

	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "alchemy_getTokenBalances",
		"headers": map[string]string{
			"Content-Type": "application/json",
		},
		"params": []interface{}{
			address,
			"erc20",
		},
		"id": 42,
	}

	payload, _ := json.Marshal(data)

	request, _ := http.NewRequest("POST", baseURL, bytes.NewBuffer(payload))
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.WithError(err).Error("cannot connect to server")
		return
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != 200 {
		log.WithError(err).Error("received unexpected response status code from alchemy_getTokenBalances request", response.StatusCode)
	}

	responseData, _ := io.ReadAll(response.Body)
	var result map[string]interface{}
	_ = json.Unmarshal(responseData, &result)

	var tokenBalances GetTokenBalanceResponse
	_ = json.Unmarshal(responseData, &tokenBalances)

	var (
		concurrentGoroutines = 5
		semaphore            = make(chan struct{}, concurrentGoroutines)
		wg                   sync.WaitGroup
		mu                   sync.Mutex
	)

	var ftbs []*FullTokenBalance

	for i := range tokenBalances.Result.TokenBalances {
		if tokenBalances.Result.TokenBalances[i].TokenBalance != "0x0000000000000000000000000000000000000000000000000000000000000000" {
			token := tokenBalances.Result.TokenBalances[i]

			// Acquisisci un semaforo prima di avviare una nuova goroutine
			semaphore <- struct{}{}
			wg.Add(1)
			go func(token TokenBalance) {
				defer wg.Done()
				defer func() { <-semaphore }() // Rilascia il semaforo al termine
				ftb, err := getInfo(browser, token)
				if err != nil {
					log.WithError(err).Error("failed to get token info")
				} else {
					mu.Lock()
					ftbs = append(ftbs, ftb)
					mu.Unlock()
				}
			}(token)
		}
	}
	wg.Wait()

	// ftbs è pronto
	return nil, nil
}

func getInfo(browser playwright.Browser, balance TokenBalance) (*FullTokenBalance, error) {

	page, err := browser.NewPage()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = page.Close()
	}()
	page.SetDefaultTimeout(1000 * 40)
	_, err = page.Goto(fmt.Sprintf("https://etherscan.io/token/%s", balance.ContractAddress))
	if err != nil {
		return nil, err
	}
	decimalLocator := page.Locator("#ContentPlaceHolder1_divSummary > div.row.g-3.mb-4 > div:nth-child(3) > div > div > div:nth-child(2) > h4 > b")
	decimalString, err := decimalLocator.TextContent()
	if err != nil {
		return nil, err
	}

	tokenNameLocator := page.Locator("xpath=/html/body/main/section[1]/div/div[1]/div/span[1]")
	tokenName, err := tokenNameLocator.TextContent()
	if err != nil {
		return nil, err
	}
	util.CleanText(&tokenName)

	decimal, err := strconv.Atoi(decimalString)
	if err != nil {
		return nil, err
	}
	intValue := new(big.Int)
	intValue.SetString(balance.TokenBalance, 0)
	scale := new(big.Float).SetInt(big.NewInt(10).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil))
	floatValue := new(big.Float).SetInt(intValue)
	floatValue.Quo(floatValue, scale)
	approxValue := floatValue.Text('f', decimal)

	ftb := FullTokenBalance{
		Contract: balance.ContractAddress,
		Name:     tokenName,
		Amount:   func() float64 { s, _ := strconv.ParseFloat(approxValue, 64); return s }(),
	}
	return &ftb, nil
}

type FullTokenBalance struct {
	Contract string  `json:"contract"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
}
