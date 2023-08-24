package main

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"strings"
)

type Transaction struct {
	TransactionHash    string
	TxType             string
	Method             string
	Block              string
	AgeMillisA         string
	AgeTimestamp       string
	AgeDistanceFromNow string
	GasPrice           string
	From               string
	To                 string
	Amount             string
	Value              string
	Asset              string
	TxnFee             string
}

func main() {
	pw, err := playwright.Run()
	if err != nil {
		fmt.Println("Errore durante l'avvio di Playwright:", err)
		return
	}
	defer pw.Stop()

	if err = playwright.Install(); err != nil {
		fmt.Println("Errore durante l'installazione di Playwright:", err)
		return
	}

	browser, err := pw.Firefox.Launch()
	if err != nil {
		fmt.Println("Errore durante l'avvio del browser:", err)
		return
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		fmt.Println("Errore durante la creazione della pagina:", err)
		return
	}

	page.SetDefaultTimeout(1000 * 40)

	_, err = page.Goto("https://etherscan.io/advanced-filter?fadd=0x905615DE62BE9B1a6582843E8ceDeDB6BDA42367&tadd=0x905615DE62BE9B1a6582843E8ceDeDB6BDA42367&txntype=2")
	if err != nil {
		fmt.Println("Errore durante la navigazione alla pagina:", err)
		return
	}

	_, err = page.WaitForSelector("table.table tbody tr")
	if err != nil {
		fmt.Println("Errore durante l'attesa della tabella:", err)
		return
	}

	rows, err := page.QuerySelectorAll("table.table tbody tr")
	if err != nil {
		fmt.Println("Errore durante il recupero delle righe:", err)
		return
	}

	var transactions []Transaction

	for _, row := range rows {
		cells, err := row.QuerySelectorAll("td")
		if err != nil {
			fmt.Println("Errore durante il recupero delle celle:", err)
			continue
		}

		var rowData Transaction
		var cellData []string
		for _, cell := range cells {
			text, err := cell.TextContent()
			if err != nil {
				fmt.Println("Errore durante la lettura delle celle:", err)
			}
			text = strings.TrimSpace(text)
			text = strings.Join(strings.Fields(text), " ")
			if len(text) == 0 {
				continue
			}
			cellData = append(cellData, text)
		}

		if len(cellData) != 14 {
			fmt.Println("Errore: numero di colonne non valido")
			continue
		}

		rowData = Transaction{
			TransactionHash:    cellData[0],
			TxType:             cellData[1],
			Method:             cellData[2],
			Block:              cellData[3],
			AgeTimestamp:       cellData[4],
			AgeDistanceFromNow: cellData[5],
			AgeMillisA:         cellData[6],
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

	for _, t := range transactions {
		fmt.Println(t)
	}
}
