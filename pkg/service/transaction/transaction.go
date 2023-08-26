package transaction

import (
	"errors"
	database2 "ethereum-mirror/pkg/database"
	"ethereum-mirror/pkg/model"
	"fmt"
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"slices"
	"sort"
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

	_, err = page.Goto(fmt.Sprintf("https://etherscan.io/advanced-filter?fadd=%s&tadd=%s&txntype=0", address, address))
	if err != nil {
		return nil, err
	}

	_, err = page.WaitForSelector("table.table tbody tr")
	if err != nil {
		return nil, err
	}

	rows, err := page.QuerySelectorAll("table.table tbody tr")
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		cells, err := row.QuerySelectorAll("td")
		if err != nil {
			fmt.Println("error retrieving cells:", err)
			continue
		}

		var rowData model.Transaction
		var cellData []string
		for _, cell := range cells {
			text, err := cell.TextContent()
			if err != nil {
				fmt.Println("error reading cells:", err)
				continue
			}
			text = strings.TrimSpace(text)
			text = strings.Join(strings.Fields(text), " ")
			if len(text) == 0 {
				continue
			}
			cellData = append(cellData, text)
		}

		if len(cellData) != 14 {
			fmt.Println("error: invalid number of columns")
			continue
		}

		rowData = model.Transaction{
			TransactionHash: cellData[0],
			TxType:          cellData[1],
			Method:          cellData[2],
			Block:           cellData[3],
			AgeTimestamp: func() time.Time {
				t, _ := time.Parse(time.DateTime, cellData[4])
				return t
			}(),
			AgeDistanceFromQuery: cellData[5],
			AgeMillis:            cellData[6],
			From:                 cellData[7],
			To:                   cellData[8],
			Amount:               cellData[9],
			Value:                cellData[10],
			Asset:                cellData[11],
			TxnFee:               cellData[12],
			GasPrice:             cellData[13],
		}

		transactions = append(transactions, rowData)
	}

	// Sort transactions by AgeTimestamp in descending order
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].AgeTimestamp.After(transactions[j].AgeTimestamp)
	})

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
		tx.TxType = transactions[i].TxType
		tx.Method = transactions[i].Method
		tx.Block = transactions[i].Block
		tx.AgeMillis = transactions[i].AgeMillis
		tx.AgeTimestamp = transactions[i].AgeTimestamp
		tx.AgeDistanceFromQuery = transactions[i].AgeDistanceFromQuery
		tx.GasPrice = transactions[i].GasPrice
		tx.From = transactions[i].From
		tx.To = transactions[i].To
		tx.Amount = transactions[i].Amount
		tx.Value = transactions[i].Value
		tx.Asset = transactions[i].Asset
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
