package transaction

import (
	"errors"
	database2 "ethereum-mirror/pkg/database"
	"ethereum-mirror/pkg/model"
	"ethereum-mirror/pkg/util"
	"fmt"
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"slices"
	"strings"
	"time"
)

// GetByLimitAndAddress returns a list of all transactions of ethereum address
// TODO
// Table has 25 rows (if limit is greater than 25, process will execute more than one page, clicking on next table page)
func GetByLimitAndAddress(browser playwright.Browser, address string) (transactions []model.Transaction, err error) {

	page, err := browser.NewPage()
	if err != nil {
		return nil, err
	}
	page.SetDefaultTimeout(1000 * 40)
	defer func() {
		_ = page.Close()
	}()

	_, err = page.Goto(fmt.Sprintf("https://etherscan.io/txs?a=%s", address))
	if err != nil {
		return nil, err
	}

	_, err = page.WaitForSelector(string(util.HeaderPageSize))
	if err != nil {
		return nil, err
	}

	pageSize, err := page.QuerySelector(string(util.HeaderPageSize))
	if err != nil {
		return nil, err
	}

	if _, err = pageSize.SelectOption(playwright.SelectOptionValues{Values: &[]string{"25"}}); err != nil {
		return nil, err
	}

	_, err = page.WaitForSelector(string(util.HeaderTable))
	if err != nil {
		return nil, err
	}

	rows, err := page.QuerySelectorAll(string(util.HeaderTable))
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		cells, err := row.QuerySelectorAll(string(util.HeaderTableData))
		if err != nil {
			return nil, err
		}

		var rowData model.Transaction
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
			if slices.Contains([]int{7, 9}, i) {
				el, err := cells[i].QuerySelector(`xpath=div/a[@class="js-clipboard link-secondary "]`)
				if err != nil {
					return nil, err
				}
				text, err = el.GetAttribute("data-clipboard-text")
				if err != nil {
					return nil, err
				}
				util.CleanText(&text)
			}
			cellData = append(cellData, text)
		}

		if len(cellData) != 12 {
			return nil, errors.New("invalid cell data columns")
		}

		rowData = model.Transaction{
			TransactionHash: cellData[0],
			Method:          cellData[1],
			Block:           cellData[2],

			AgeTimestamp: func() time.Time {
				t, _ := time.Parse(time.DateTime, cellData[3])
				return t
			}(),
			AgeDistanceFromQuery: cellData[4],
			AgeMillis:            cellData[5],
			From:                 cellData[6],
			InOut:                cellData[7],
			To:                   cellData[8],
			Value:                cellData[9],
			TxnFee:               cellData[10],
			GasPrice:             cellData[11],
		}
		transactions = append(transactions, rowData)
	}

	// Sort transactions by AgeTimestamp in descending order
	//sort.Slice(transactions, func(i, j int) bool {
	//	return transactions[i].AgeTimestamp.After(transactions[j].AgeTimestamp)
	//})

	return transactions, nil
}

// SaveNew saves only new transaction data in the database Transaction table
// transactions input must be sorted by AgeTimestamp from recent to oldest
func SaveNew(database *gorm.DB, transactions []model.Transaction) ([]database2.Transaction, error) {
	scraping := database2.Scraping{}
	if err := database.Create(&scraping).Error; err != nil {
		return nil, err
	}

	var recentTransaction database2.Transaction
	from := 0
	to := len(transactions)
	if err := database.Order("AgeTimestamp DESC").First(&recentTransaction).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	} else {
		// if it's present and its position is different from the first in the list, then take all elements until the block element
		idx := slices.IndexFunc(transactions, func(t model.Transaction) bool {
			return t.TransactionHash == recentTransaction.TransactionHash
		})
		switch idx {
		case -1:
			break
		case 0:
			return nil, nil
		}
		if idx > 0 {
			to = idx
		}
	}

	var dbTransactions []database2.Transaction
	semaphore := make(chan struct{}, 10) // Create a semaphore with a capacity of 10
	tx := database2.Transaction{
		ScrapingId: scraping.ScrapingId,
	}
	for i := from; i < to; i++ {
		// Acquire a permit from the semaphore
		semaphore <- struct{}{}

		tx.TransactionHash = transactions[i].TransactionHash
		tx.Method = transactions[i].Method
		tx.Block = transactions[i].Block
		tx.AgeMillis = transactions[i].AgeMillis
		tx.AgeTimestamp = transactions[i].AgeTimestamp
		tx.AgeDistanceFromQuery = transactions[i].AgeDistanceFromQuery
		tx.GasPrice = transactions[i].GasPrice
		tx.From = transactions[i].From
		tx.To = transactions[i].To
		tx.InOut = transactions[i].InOut
		tx.Value = transactions[i].Value
		tx.TxnFee = transactions[i].TxnFee

		go func(tx database2.Transaction) {
			defer func() {
				// Release the permit back to the semaphore
				<-semaphore
			}()

			if err := database.Create(&tx).Error; err != nil {
				log.WithError(err).Errorf("failed to save transaction %s", tx.TransactionHash)
			} else {
				dbTransactions = append(dbTransactions, tx)
			}
		}(tx)
	}

	// Wait for all goroutines to finish
	for i := 0; i < cap(semaphore); i++ {
		semaphore <- struct{}{}
	}

	return dbTransactions, nil
}
