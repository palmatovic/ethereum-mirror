package transaction_detail

import (
	"encoding/json"
	"ethereum-mirror/pkg/database"
	"ethereum-mirror/pkg/model"
	"ethereum-mirror/pkg/util"
	"fmt"
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
)

func GetByTransaction(browser playwright.Browser, transactions []database.Transaction) ([]model.TransactionDetail, error) {
	var details []model.TransactionDetail

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)
	for i := range transactions {
		wg.Add(1)
		sem <- struct{}{} // acquire a semaphore slot
		go func(transaction database.Transaction) {
			defer wg.Done()
			page, err := browser.NewPage()
			if err != nil {
				panic(err)
			}
			defer func() {
				_ = page.Close()
			}()

			page.SetDefaultTimeout(1000 * 40)

			var detail model.TransactionDetail

			_, err = page.Goto(fmt.Sprintf("https://etherscan.io/tx/%s", transaction.TransactionHash))
			if err != nil {
				panic(err)
			}
			txHash, err := util.GetObjectByPage(page, util.TxHash)
			if err != nil {
				panic(err)
			}
			detail.TransactionHash = txHash

			status, err := util.GetObjectByPage(page, util.TxStatus)
			if err != nil {
				panic(err)
			}
			detail.TransactionStatus = status

			block, err := util.GetObjectByPage(page, util.TxBlock)
			if err != nil {
				panic(err)
			}
			detail.TransactionBlock = block

			timestamp, err := util.GetObjectByPage(page, util.TxTimestamp)
			if err != nil {
				panic(err)
			}
			detail.TransactionTimestamp = timestamp

			if strings.ToLower(status) == "success" {
				action, err := util.GetObjectByPage(page, util.TxAction)
				if err != nil {
					panic(err)
				}
				detail.TransactionAction = action
			}

			from, err := util.GetObjectByPage(page, util.TxFrom)
			if err != nil {
				panic(err)
			}
			detail.TransactionFrom = from

			var to string
			to, err = util.GetObjectByPage(page, util.TxInteractedWithToSuccess)
			if err != nil {
				var errFail error
				to, errFail = util.GetObjectByPage(page, util.TxInteractedWithToFail)
				if errFail != nil {
					log.WithError(err).Errorf("cannot get To for transaction %v", transaction.TransactionHash)
					panic(err)
				}
			}
			detail.TransactionTo = to

			details = append(details, detail)

			<-sem
		}(transactions[i])
	}
	wg.Wait()

	b, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	file, err := os.Create("detail.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	_, err = file.Write(b)

	return details, nil
}
