package transaction_detail

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"transaction-extractor/pkg/database/transaction"
	"transaction-extractor/pkg/model/transaction_detail"
	"transaction-extractor/pkg/util"
)

func GetByTransaction(browser playwright.Browser, transactions []transaction.Transaction) ([]transaction_detail.TransactionDetail, error) {
	var details []transaction_detail.TransactionDetail

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)
	for i := range transactions {
		wg.Add(1)
		sem <- struct{}{} // acquire a semaphore slot
		go func(transaction transaction.Transaction) {
			defer wg.Done()
			page, err := browser.NewPage()
			if err != nil {
				panic(err)
			}
			defer func() {
				_ = page.Close()
			}()

			page.SetDefaultTimeout(1000 * 40)

			var detail transaction_detail.TransactionDetail

			_, err = page.Goto(fmt.Sprintf("https://etherscan.io/tx/%s", transaction.TxHash))
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
					log.WithError(err).Panicf("error getting transaction action for transaction %v", transaction.TxHash)
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
					log.WithError(err).Errorf("cannot get To for transaction %v", transaction.TxHash)
					panic(err)
				}
			}
			detail.TransactionTo = to

			details = append(details, detail)

			<-sem
		}(transactions[i])
	}
	wg.Wait()

	return details, nil
}
