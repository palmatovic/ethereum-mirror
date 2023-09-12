package address_status

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	address_status_db "transaction-extractor/pkg/database/address_status"
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
func GetAddressStatus(db *gorm.DB, browser playwright.Browser, address string, alchemyApiKey string) (addressStatuses []address_status.AddressStatus, err error) {
	// Alchemy URL
	var baseURL = fmt.Sprintf("https://eth-mainnet.g.alchemy.com/v2/%s", alchemyApiKey)

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
		logrus.WithError(err).Error("cannot connect to server")
		return
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != 200 {
		logrus.WithError(err).Error("received unexpected response status code from alchemy_getTokenBalances request", response.StatusCode)
	}

	responseData, _ := io.ReadAll(response.Body)

	var tokenBalances GetTokenBalanceResponse
	_ = json.Unmarshal(responseData, &tokenBalances)

	var (
		concurrentGoroutines = 5
		semaphore            = make(chan struct{}, concurrentGoroutines)
		wg                   sync.WaitGroup
		mu                   sync.Mutex
	)

	for i := range tokenBalances.Result.TokenBalances {
		if tokenBalances.Result.TokenBalances[i].TokenBalance != "0x0000000000000000000000000000000000000000000000000000000000000000" {
			token := tokenBalances.Result.TokenBalances[i]
			semaphore <- struct{}{}
			wg.Add(1)
			go func(t TokenBalance) {
				defer wg.Done()
				defer func() { <-semaphore }()
				var addressStatusDb address_status_db.AddressStatus
				err = db.Where("AddressId = ? AND TokenContractAddress = ?", address, token.ContractAddress).First(&addressStatusDb).Error
				if err != nil {
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						return
					}
				}
				if addressStatusDb.TokenAmountHex != t.TokenBalance {
					addressStatus, err := getInfo(browser, t)
					if err != nil {
						logrus.WithError(err).Error("failed to get token info")
					} else {
						mu.Lock()
						addressStatuses = append(addressStatuses, *addressStatus)
						mu.Unlock()
					}
				}
			}(token)
		}
	}
	wg.Wait()

	return addressStatuses, nil
}

func getInfo(browser playwright.Browser, balance TokenBalance) (*address_status.AddressStatus, error) {

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

	ftb := address_status.AddressStatus{
		Contract:  balance.ContractAddress,
		Name:      tokenName,
		Amount:    func() float64 { s, _ := strconv.ParseFloat(approxValue, 64); return s }(),
		AmountHex: balance.TokenBalance,
		Symbol: func() string {
			if strings.Contains(tokenName, "(") && strings.Contains(tokenName, ")") {
				start := strings.Index(tokenName, "(")
				end := strings.Index(tokenName, ")")
				content := tokenName[start+1 : end]
				return content
			} else {
				return tokenName
			}
		}(),
	}
	return &ftb, nil
}

func UpsertAddressStatus(db *gorm.DB, address string, addressStatuses []address_status.AddressStatus) ([]address_status_db.AddressStatus, error) {
	var addressStatusesDb []address_status_db.AddressStatus
	for i := range addressStatuses {
		addressStatusesDb = append(addressStatusesDb, address_status_db.AddressStatus{
			AddressId:            address,
			TokenContractAddress: addressStatuses[i].Contract,
			TokenName:            addressStatuses[i].Name,
			TokenSymbol:          addressStatuses[i].Symbol,
			TokenAmount:          addressStatuses[i].Amount,
			TokenAmountHex:       addressStatuses[i].AmountHex,
		})
	}
	for j := range addressStatusesDb {
		var asd address_status_db.AddressStatus
		var err error
		err = db.Where("AddressId = ? AND TokenContractAddress = ?", address, addressStatusesDb[j].TokenContractAddress).First(&asd).Error
		if err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				if err = db.Create(&addressStatusesDb[j]).Error; err != nil {
					return nil, err
				}
				// inserire nella la tabella delle transazioni che hanno portato a quel token amount x il token contract address e l'addressId (etherscan puoi filtrare per contratto e holder wallet) (che non esiste)

			} else {
				return nil, err
			}
		} else {
			if asd.TokenAmount != addressStatusesDb[j].TokenAmount {
				// aggiorna la tabelle address_status con il nuovo token amont e aggiorna la tabella delle transazioni con solo le transazioni nuove (che non esiste vedi sopra)
			}
		}
	}
	return addressStatusesDb, nil
}
