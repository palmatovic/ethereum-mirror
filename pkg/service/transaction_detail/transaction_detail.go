package transaction_detail

import (
	"ethereum-mirror/pkg/database"
	"ethereum-mirror/pkg/model"
	"ethereum-mirror/pkg/util"
	"fmt"
	"github.com/playwright-community/playwright-go"
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
			defer page.Close()
			page.SetDefaultTimeout(1000 * 40)

			var detail model.TransactionDetail
			//log.Infof("retrieving detail for %s", transaction.TransactionHash)

			_, err = page.Goto(fmt.Sprintf("https://etherscan.io/tx/%s", transaction.TransactionHash))
			if err != nil {
				panic(err)
			}
			txHash, err := util.GetObjectByPage(page, util.TxHash)
			if err != nil {
				panic(err)
			}
			detail.TransactionHash = txHash

			action, err := util.GetObjectByPage(page, util.TxAction)
			if err != nil {
				panic(err)
			}
			detail.TransactionAction = action

			status, err := util.GetObjectByPage(page, util.TxStatus)
			if err != nil {
				panic(err)
			}
			detail.Status = status

			details = append(details, detail)
			<-sem
		}(transactions[i])
	}
	wg.Wait()

	fmt.Println(details)
	return details, nil
}
