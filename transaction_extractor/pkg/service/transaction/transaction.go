package transaction

import (
	"errors"
	"fmt"
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"slices"
	"strings"
	"time"
	database2 "transaction-extractor/pkg/database/scraping"
	transaction2 "transaction-extractor/pkg/database/transaction"
	"transaction-extractor/pkg/model/transaction"
	"transaction-extractor/pkg/util"
)

// GetByAddress returns a list of all transactions of ethereum address
func GetByAddress(browser playwright.Browser, address string) (transactions []transaction.Transaction, err error) {

	page, err := browser.NewPage()
	if err != nil {
		return nil, err
	}
	page.SetDefaultTimeout(1000 * 40)
	defer func() {
		_ = page.Close()
	}()

	_, err = page.Goto(fmt.Sprintf("https://etherscan.io/advanced-filter?fadd=%s&tadd=%s&p=3&fs=1", address, address))
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

	if _, err = pageSize.SelectOption(playwright.SelectOptionValues{Values: &[]string{"100"}}); err != nil {
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

		var rowData transaction.Transaction
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

		if len(cellData) != 14 {
			return nil, errors.New("invalid cell data columns")
		}

		if slices.IndexFunc(transactions, func(t transaction.Transaction) bool {
			return t.TxHash == cellData[0]
		}) == -1 {
			rowData = transaction.Transaction{
				TxHash: cellData[0],
				TxType: cellData[1],
				Method: cellData[2],
				Block:  cellData[3],
				AgeTimestamp: func() time.Time {
					t, _ := time.Parse(time.DateTime, cellData[4])
					return t
				}(),
				AgeDistanceFromQuery: cellData[5],
				AgeMillis:            cellData[6],
				From:                 cellData[7],
				//InOut:                cellData[7],
				To:            cellData[8],
				Amount:        cellData[9],
				Value:         cellData[10],
				Asset:         cellData[11],
				TxnFee:        cellData[12],
				GasPrice:      cellData[13],
				WalletAddress: address,
			}
			transactions = append(transactions, rowData)
		}
	}
	return transactions, nil
}

// SaveNew saves only new transaction data in the database Transaction table
// transactions input must be sorted by AgeTimestamp from recent to oldest
func SaveNew(database *gorm.DB, transactions []transaction.Transaction) ([]transaction2.Transaction, error) {

	scraping := database2.Scraping{}
	if err := database.Create(&scraping).Error; err != nil {
		return nil, err
	}

	var recentTransaction transaction2.Transaction
	from := 0
	to := len(transactions)
	if err := database.Order("AgeTimestamp DESC").First(&recentTransaction).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	} else {
		// if it's present and its position is different from the first in the list, then take all elements until the block element
		idx := slices.IndexFunc(transactions, func(t transaction.Transaction) bool {
			return t.TxHash == recentTransaction.TxHash
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

	var dbTransactions []transaction2.Transaction
	semaphore := make(chan struct{}, 10) // Create a semaphore with a capacity of 10

	for i := from; i < to; i++ {
		// Acquire a permit from the semaphore
		semaphore <- struct{}{}
		tx := transaction2.Transaction{
			TxHash:               transactions[i].TxHash,
			TxType:               transactions[i].TxType,
			Method:               transactions[i].Method,
			Block:                transactions[i].Block,
			AgeTimestamp:         transactions[i].AgeTimestamp,
			AgeDistanceFromQuery: transactions[i].AgeDistanceFromQuery,
			AgeMillis:            transactions[i].AgeMillis,
			From:                 transactions[i].From,
			To:                   transactions[i].To,
			Amount:               transactions[i].Amount,
			Value:                transactions[i].Value,
			Asset:                transactions[i].Asset,
			TxnFee:               transactions[i].TxnFee,
			GasPrice:             transactions[i].GasPrice,
			ScrapingId:           scraping.ScrapingId,
			WalletAddress:        transactions[i].WalletAddress,
		}
		go func(tx transaction2.Transaction) {
			defer func() {
				// Release the permit back to the semaphore
				<-semaphore
			}()

			if err := database.Create(&tx).Error; err != nil {
				log.WithError(err).Errorf("failed to save transaction %s", tx.TxHash)
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
