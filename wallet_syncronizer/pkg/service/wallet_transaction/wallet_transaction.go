package wallet_transaction

import (
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"time"
	"wallet-syncronizer/pkg/database/wallet_token"
	"wallet-syncronizer/pkg/database/wallet_transaction"
	string2 "wallet-syncronizer/pkg/util/string"
)

func FindOrCreateWalletTransactions(db *gorm.DB, walletTokens []wallet_token.WalletToken, browser playwright.Browser) (err error) {

	var (
	//concurrentGoroutines = 10
	//semaphore            = make(chan struct{}, concurrentGoroutines)
	//wg                   sync.WaitGroup
	//mutex                sync.Mutex
	)
	var ats []wallet_transaction.WalletTransaction

	for _, walletToken := range walletTokens {
		//semaphore <- struct{}{}
		//wg.Add(1)
		//go func(walletToken wallet_token.WalletToken) {
		var page playwright.Page
		var bctx playwright.BrowserContext
		bctx, err = browser.NewContext()
		page, err = bctx.NewPage()
		//defer func() {
		//	//wg.Done()
		//	//<-semaphore
		//}()
		if err != nil {
			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot open page")
			_ = page.Close()
			return
		}

		_, err = page.Goto("https://www.defined.fi/eth/" + walletToken.TokenId)
		if err != nil {
			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot open page")
			_ = page.Close()
			return
		}

		err = page.WaitForLoadState()
		if err != nil {
			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot wait for page")
			_ = page.Close()
			return
		}
		page.WaitForTimeout(2000)
		//page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[3]/div/div/div[2]/div[1]/div/div")
		err = page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[2]/div[5]").WaitFor(playwright.LocatorWaitForOptions{
			Timeout: playwright.Float(2000),
		})
		if err != nil {
			// TODO:
			notFoundLocator := page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/h2")
			notFound, err2 := notFoundLocator.TextContent()
			if err2 != nil {
				return err2
			}
			string2.CleanText(&notFound)
			if strings.ToLower(notFound) == "not found" {
				// aggiungi come scam token, anche se in realtà andrebbe filtrato a priori una volta ottenuti i dati da go plus
				// nel service token
				_ = page.Close()
				continue
			} else {
				return err
			}
		}
		page.WaitForTimeout(2000)
		expandTable := page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[2]/div[5]")
		err = expandTable.Click()
		if err != nil {
			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot click on expandTable locator")
			_ = page.Close()
			//// TODO:
			//notFoundLocator := page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/h2")
			//notFound, err2 := notFoundLocator.TextContent()
			//if err2 != nil {
			//	return err2
			//}
			//string.CleanText(&notFound)
			//if strings.ToLower(notFound) == "not found" {
			//	// aggiungi come scam token, anche se in realtà andrebbe filtrato a priori una volta ottenuti i dati da go plus
			//	// nel service token
			//	continue
			//} else {
			//	return err
			//}
			return err
		}

		_ = page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[3]/div/div/div[1]/div[7]/span/button").WaitFor(playwright.LocatorWaitForOptions{
			Timeout: playwright.Float(2000),
		})
		filterButton := page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[3]/div/div/div[1]/div[7]/span/button")
		err = filterButton.Click()
		if err != nil {
			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot click on filter locator")
			_ = page.Close()
			return
		}

		_ = page.Locator("xpath=//html/body/div[7]/div[3]/form/div[1]/div/input").WaitFor(playwright.LocatorWaitForOptions{
			Timeout: playwright.Float(2000),
		})

		inputFiter := page.Locator("xpath=//html/body/div[7]/div[3]/form/div[1]/div/input")
		err = inputFiter.Fill(walletToken.WalletId)
		if err != nil {
			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot fill filter locator")
			_ = page.Close()
			return err
		}

		applyFilter := page.Locator("xpath=//html/body/div[7]/div[3]/form/div[2]/button[2]")
		err = applyFilter.Click()
		if err != nil {
			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot apply filter locator")
			_ = page.Close()
			return err
		}

		tableXpath := "xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[3]/div/div/div[2]/div[1]/div/div"

		_ = page.Locator(tableXpath).WaitFor(playwright.LocatorWaitForOptions{
			Timeout: playwright.Float(2000),
		})

		table := page.Locator(tableXpath)
		rows, err := table.Locator("xpath=/div").All()
		if err != nil {
			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot get all rows locator")
			_ = page.Close()
			return err
		}

		for _, r := range rows {
			cols, err := r.Locator("xpath=/div/div").All()
			if err != nil {
				logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot get columns row %d locator", r)
				return err
			}

			var at wallet_transaction.WalletTransaction
			at.WalletId = walletToken.WalletId
			at.Asset = walletToken.TokenId

			for colNum, _ := range cols {

				if colNum == 0 {
					at.TxType, err = cols[colNum].InnerText()
					if err != nil {
						logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot get inner text row %d column %d", r, colNum)
						return err
					}

				}

				if colNum == 1 {
					var strPrice string
					strPrice, err = cols[colNum].TextContent()
					if err != nil {
						logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot get text content row %d column %d", r, colNum)
						return err
					}
					strPrice = strings.ReplaceAll(strPrice, ",", "")
					strPrice = strings.ReplaceAll(strPrice, "$", "")
					string2.CleanText(&strPrice)
					at.Price, err = string2.ParseScript(strPrice)
					if err != nil {
						logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot parse price: row %d column %d", r, colNum)
						return err
					}
					//if at.Price, err = strconv.ParseFloat(strPrice, 64); err != nil {
					//	at.Price, err = string2.ParseScript(strPrice)
					//	if err != nil {
					//		string2.CleanTextWithRemoveUnicodeSpaces(&strPrice)
					//		at.Price, err = string2.ParseScript(strPrice)
					//		if err != nil {
					//			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot parse price: row %d column %d", r, colNum)
					//			return err
					//		}
					//	}
					//}

				}

				if colNum == 2 {
					var strAmount string
					strAmount, err = cols[colNum].TextContent()
					if err != nil {
						logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot get text content row %d column %d", r, colNum)
						return err
					}
					strAmount = strings.ReplaceAll(strAmount, ",", "")
					strAmount = strings.ReplaceAll(strAmount, "$", "")
					string2.CleanText(&strAmount)
					at.Amount, err = string2.ParseScript(strAmount)
					if err != nil {
						logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot parse amount: row %d column %d", r, colNum)
						return err
					}
					//if at.Amount, err = strconv.ParseFloat(strAmount, 64); err != nil {
					//	at.Amount, err = string2.ParseScript(strAmount)
					//	if err != nil {
					//		string2.CleanTextWithRemoveUnicodeSpaces(&strAmount)
					//		at.Amount, err = string2.ParseScript(strAmount)
					//		if err != nil {
					//			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot parse amount: row %d column %d", r, colNum)
					//			return err
					//		}
					//	}
					//}
				}

				if colNum == 3 {
					var strTotal string
					strTotal, err = cols[colNum].TextContent()
					if err != nil {
						logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot get text content row %d column %d", r, colNum)
						return err
					}
					strTotal = strings.ReplaceAll(strTotal, ",", "")
					strTotal = strings.ReplaceAll(strTotal, "$", "")
					string2.CleanText(&strTotal)
					at.Total, err = string2.ParseScript(strTotal)
					if err != nil {
						logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot parse total: row %d column %d", r, colNum)
						return err
					}
					//if at.Total, err = strconv.ParseFloat(strTotal, 64); err != nil {
					//	at.Total, err = string2.ParseScript(strTotal)
					//	if err != nil {
					//		string2.CleanTextWithRemoveUnicodeSpaces(&strTotal)
					//		at.Total, err = string2.ParseScript(strTotal)
					//		if err != nil {
					//			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot parse total: row %d column %d", r, colNum)
					//			return err
					//		}
					//	}
					//}
				}

				if colNum == 4 {
					continue
				}

				if colNum == 5 {

					var colSpan = cols[colNum].Locator("xpath=/span")

					ageTimestamp, err := colSpan.GetAttribute("aria-label")
					if err != nil {
						logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot get attribute aria-label: row %d column %d", r, colNum)
						return err
					}

					at.AgeTimestamp, err = time.Parse(time.DateTime, ageTimestamp)
					if err != nil {
						logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot parse age timestamp: row %d column %d", r, colNum)
						return err
					}

				}

				//if colNum == 7 {
				//
				//	var colA = cols[colNum].Locator("xpath=/span/a")
				//
				//	err := colA.Click()
				//	if err != nil {
				//		return err
				//	}
				//
				//	newPageInterface, err := browser.Contexts()[0].WaitForEvent("page")
				//	if err != nil {
				//		return err
				//	}
				//	np := newPageInterface.(playwright.Page)
				//
				//	xpathActionsContainer := "//div[@id='wrapperContent']"
				//	xpathActions := "xpath=/div/div"
				//	//np = browser.Contexts()[0].Pages()[0]
				//
				//	ac := np.Locator(xpathActionsContainer)
				//	allActions, err := ac.Locator(xpathActions).All()
				//	if err != nil {
				//		return err
				//	}
				//
				//	var steps [][]string
				//	for _, aa := range allActions {
				//		var step []string
				//		txts, err := aa.InnerText()
				//		if err != nil {
				//			return err
				//		}
				//		fromString := strings.Split(txts, "\n")[0]
				//
				//		step = append(step, txts)
				//		steps = append(steps, step)
				//	}
				//
				//}

			}

			ats = append(ats, at)
			logrus.Info(at)
			//logrus.Info(string(func() []byte {
			//	byteArray, _ := json.Marshal(at)
			//	return byteArray
			//}()))
		}

		//logrus.Info(len(ats))

		//}(walletToken)
		_ = page.Close()
	}
	//wg.Wait()
	//return db.Create(&ats).Error
	return nil
}
