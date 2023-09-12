package wallet_balance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Balance struct {
	TokenContractAddress string `json:"contractAddress"`
	TokenBalance         string `json:"tokenBalance"`
}

type WalletBalanceResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  Result `json:"result"`
}

type Result struct {
	Address        string    `json:"wallet"`
	WalletBalances []Balance `json:"tokenBalances"`
}

func GetWalletBalances(wallet string, apiKey string) (walletBalance WalletBalanceResponse, err error) {
	// Alchemy URL
	var baseURL = fmt.Sprintf("https://eth-mainnet.g.alchemy.com/v2/%s", apiKey)

	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "alchemy_getTokenBalances",
		"headers": map[string]string{
			"Content-Type": "application/json",
		},
		"params": []interface{}{
			wallet,
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
		return walletBalance, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != 200 {
		logrus.WithError(err).Error("received unexpected response status code from alchemy_getTokenBalances request", response.StatusCode)
	}

	responseData, _ := io.ReadAll(response.Body)

	_ = json.Unmarshal(responseData, &walletBalance)
	return walletBalance, nil
}
