package api

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

type Env struct {
}

func (e *Env) Start() (response interface{}, err error) {

	url := "https://etherscan.io/advanced-filter?fadd=0x905615DE62BE9B1a6582843E8ceDeDB6BDA42367&tadd=0x905615DE62BE9B1a6582843E8ceDeDB6BDA42367&txntype=2"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	tableCardSelection := doc.Find("#ContentPlaceHolder1_tableCard")

	if tableCardSelection.Length() > 0 {
		tableSelection := tableCardSelection.Find(".table.table-hover.mb-0")

		if tableSelection.Length() > 0 {
			fmt.Println("Tabella trovata:", tableSelection)
		} else {
			fmt.Println("Nessuna tabella con classi 'table table-hover mb-0' trovata all'interno di ContentPlaceHolder1_tableCard.")
		}
	} else {
		fmt.Println("Nessun elemento con ID ContentPlaceHolder1_tableCard trovato.")
	}
	return nil, nil
}
