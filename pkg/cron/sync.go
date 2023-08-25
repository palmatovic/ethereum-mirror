package cron

import (
	"errors"
	database2 "ethereum-mirror/pkg/database"
	"ethereum-mirror/pkg/model"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"gorm.io/gorm"
	"slices"
	"strings"
	"time"
)

type Env struct {
	playwright.Browser
	Database *gorm.DB
}

func (e *Env) SyncTransactions() (response interface{}, err error) {

	transactions, err := getLast25Transactions(e.Browser)
	if err != nil {
		return nil, err
	}
	if err = saveNewTransactions(e.Database, transactions); err != nil {
		return nil, err
	}

	return nil, nil
}

func saveNewTransactions(database *gorm.DB, transactions []model.Transaction) error {
	scraping := database2.Scraping{}
	if err := database.Create(&scraping).Error; err != nil {
		return err
	}

	var recentTransaction database2.Transaction
	from := 0
	to := len(transactions)
	if err := database.Order("AgeTimestamp DESC").First(&recentTransaction).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
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
			return nil
		}
		if idx > 0 {
			to = idx
		}
	}

	var dbTransactions []database2.Transaction
	for i := from; i < to; i++ {
		dbTransactions = append(dbTransactions, database2.Transaction{
			TransactionHash:    transactions[i].TransactionHash,
			ScrapingId:         scraping.ScrapingId,
			TxType:             transactions[i].TxType,
			Method:             transactions[i].Method,
			Block:              transactions[i].Block,
			AgeMillis:          transactions[i].AgeMillis,
			AgeTimestamp:       transactions[i].AgeTimestamp,
			AgeDistanceFromNow: transactions[i].AgeDistanceFromNow,
			GasPrice:           transactions[i].GasPrice,
			From:               transactions[i].From,
			To:                 transactions[i].To,
			Amount:             transactions[i].Amount,
			Value:              transactions[i].Value,
			Asset:              transactions[i].Asset,
			TxnFee:             transactions[i].TxnFee})
	}
	if err := database.Create(&dbTransactions).Error; err != nil {
		return err
	}

	return nil
}

func getLast25Transactions(browser playwright.Browser) ([]model.Transaction, error) {
	page, err := browser.NewPage()
	if err != nil {
		return nil, err
	}

	page.SetDefaultTimeout(1000 * 40)

	_, err = page.Goto("https://etherscan.io/advanced-filter?fadd=0x905615DE62BE9B1a6582843E8ceDeDB6BDA42367&tadd=0x905615DE62BE9B1a6582843E8ceDeDB6BDA42367&txntype=0")
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

	var transactions []model.Transaction

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
			}
			text = strings.TrimSpace(text)
			text = strings.Join(strings.Fields(text), " ")
			if len(text) == 0 {
				continue
			}
			cellData = append(cellData, text)
		}

		if len(cellData) != 14 {
			return nil, errors.New("error: invalid number of columns")
		}

		rowData = model.Transaction{
			TransactionHash:    cellData[0],
			TxType:             cellData[1],
			Method:             cellData[2],
			Block:              cellData[3],
			AgeTimestamp:       func() time.Time { t, _ := time.Parse(time.RFC3339Nano, cellData[4]); return t }(),
			AgeDistanceFromNow: cellData[5],
			AgeMillis:          cellData[6],
			From:               cellData[7],
			To:                 cellData[8],
			Amount:             cellData[9],
			Value:              cellData[10],
			Asset:              cellData[11],
			TxnFee:             cellData[12],
			GasPrice:           cellData[13],
		}

		transactions = append(transactions, rowData)
	}

	return transactions, nil
}
