package wallet_transaction

import (
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"sync"
	"time"
	"wallet-synchronizer/pkg/database/wallet_token"
	"wallet-synchronizer/pkg/database/wallet_transaction"
)

func FindOrCreateWalletTransactions(db *gorm.DB, walletTokens []wallet_token.WalletToken, browser playwright.Browser) (err error) {

	var (
		concurrentGoroutines = 5
		semaphore            = make(chan struct{}, concurrentGoroutines)
		wg                   sync.WaitGroup
	)

	for _, wt := range walletTokens {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(walletToken wallet_token.WalletToken) {
			var ats []wallet_transaction.WalletTransaction
			var page playwright.Page
			var bctx playwright.BrowserContext
			bctx, err = browser.NewContext()
			page, err = bctx.NewPage()
			defer func() {
				_ = page.Close()
				<-semaphore
				wg.Done()
			}()
			if err != nil {
				logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot open page")
				_ = page.Close()
				return
			}

			_, err = page.Goto("https://www.defined.fi/eth/" + walletToken.TokenId)
			if err != nil {
				logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot open page")
				return
			}

			err = doWait(&page, true, true, true)
			if err != nil {
				return
			}

			ats, err = FindOrCreateWalletTransactionByLiquidityPool(walletToken, page)
			if err != nil {
				logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("failed to create wallet transactions by liquidity pool")
				return
			}
			if len(ats) > 0 {
				err = db.Create(&ats).Error
				if err != nil {
					logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("failed to create wallet transactions")
					return
				}
			}

		}(wt)
	}
	wg.Wait()
	return
}

func FindOrCreateWalletTransactionByLiquidityPool(walletToken wallet_token.WalletToken, page playwright.Page) (ats []wallet_transaction.WalletTransaction, err error) {

	//fai scraping per prendere i pool di liquidità e per ogni pool avvia la funzione sotto
	page.WaitForTimeout(3000)

	_, err = locatorWithRetryCombined(&page, "xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[2]/div[5]", 3)

	//err = page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[2]/div[5]").WaitFor(playwright.LocatorWaitForOptions{
	//	State:   playwright.WaitForSelectorStateAttached,
	//	Timeout: playwright.Float(5000),
	//})
	//err = page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[2]/div[5]").WaitFor(playwright.LocatorWaitForOptions{
	//	State:   playwright.WaitForSelectorStateVisible,
	//	Timeout: playwright.Float(5000),
	//})
	if err != nil {
		return nil, err
	}

	poolSelector, err := locatorWithRetryCombined(&page, "xpath=//html/body/div[1]/div[2]/div/div[2]/div[2]/div[1]/div[2]/div[2]/p", 3)

	//err = page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[2]/div[1]/div[2]/div[2]/p").WaitFor(playwright.LocatorWaitForOptions{
	//	State:   playwright.WaitForSelectorStateAttached,
	//	Timeout: playwright.Float(5000),
	//})
	//
	//err = page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[2]/div[1]/div[2]/div[2]/p").WaitFor(playwright.LocatorWaitForOptions{
	//	State:   playwright.WaitForSelectorStateVisible,
	//	Timeout: playwright.Float(5000),
	//})

	if err != nil {
		return nil, err
	}

	//poolSelector := page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[2]/div[1]/div[2]/div[2]/p")

	err = poolSelector.Click()
	if err != nil {
		return nil, err
	}

	err = doWait(&page, true, true, true)
	if err != nil {
		return nil, err
	}
	listPoolLocator, err := page.Locator("xpath=//html/body/div[5]/div[3]/ul/li").All()
	if err != nil {
		return nil, err
	}

	for i, li := range listPoolLocator {

		if i != 0 {
			err = poolSelector.Click()
			if err != nil {
				return nil, err
			}
			err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
				State:   playwright.LoadStateLoad,
				Timeout: playwright.Float(10000),
			})
			if err != nil {
				logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot wait for page")
				return
			}
		}
		poolName, err := li.TextContent()
		if err != nil {
			return nil, err
		}

		err = li.Click()
		if err != nil {
			return nil, err
		}

		err = doWait(&page, true, true, true)
		if err != nil {
			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot wait for page")
			return nil, err
		}

		if i == 0 {
			expandTable := page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[2]/div[5]")
			err = expandTable.Click()
			if err != nil {
				logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot click on expandTable locator")
				_ = page.Close()
				return nil, err
			}
		}

		err = page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[3]/div/div/div[1]/div[7]/span/button").WaitFor(playwright.LocatorWaitForOptions{
			Timeout: playwright.Float(2000),
		})
		if err != nil {
			return nil, err
		}

		filterButton := page.Locator("xpath=//html/body/div[1]/div[2]/div/div[2]/div[1]/div[3]/div/div/div[1]/div[7]/span/button")
		err = filterButton.Click()
		if err != nil {
			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot click on filter locator")
			_ = page.Close()
			return nil, err
		}

		err = doWait(&page, false, false, true)
		if err != nil {
			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot wait for page")
			return nil, err
		}

		_ = page.Locator("xpath=//input[@placeholder='Address']").WaitFor(playwright.LocatorWaitForOptions{
			Timeout: playwright.Float(4000),
		})

		inputFiter := page.Locator("xpath=//input[@placeholder='Address']")

		err = inputFiter.Fill(walletToken.WalletId)

		if err != nil {
			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot fill filter locator")
			return nil, err
		}

		applyFilter := page.Locator("xpath=//form[@class='css-sevhfp']/div[2]/button[2]")
		err = applyFilter.Click()
		if err != nil {
			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot apply filter locator")
			return nil, err
		}
		err = doWait(&page, true, true, true)
		if err != nil {
			logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot wait for page")
			return nil, err
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
			return nil, err
		}

		for rowNum, r := range rows {

			cols, err := r.Locator("xpath=/div/div").All()
			if err != nil {
				logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot get columns row %d locator", rowNum)
				return nil, err
			}

			var at wallet_transaction.WalletTransaction
			at.Pool = strings.Split(poolName, " ")[0]
			at.WalletId = walletToken.WalletId
			at.TokenId = walletToken.TokenId

			for colNum, _ := range cols {

				if colNum == 0 {
					at.TxType, err = cols[colNum].InnerText()
					if err != nil {
						logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot get inner text row %d column %d", rowNum, colNum)
						return nil, err
					}

				}

				if colNum == 1 {
					var strPrice string
					strPrice, err = cols[colNum].TextContent()
					if err != nil {
						logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot get text content row %d column %d", rowNum, colNum)
						return nil, err
					}
					strPrice = strings.ReplaceAll(strPrice, ",", "")
					strPrice = strings.ReplaceAll(strPrice, "$", "")
					CleanText(&strPrice)
					at.Price, err = ParseScript(strPrice)
					if err != nil {
						//logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot parse price: %s, row %d column %d. trying to remove unicode spaces", strPrice, rowNum, colNum)
						CleanTextWithRemoveUnicodeSpaces(&strPrice)
						at.Price, err = ParseScript(strPrice)
						if err != nil {
							logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot parse price %s: row %d column %d", strPrice, rowNum, colNum)
							return nil, err
						}
					}

				}

				if colNum == 2 {
					var strAmount string
					strAmount, err = cols[colNum].TextContent()
					if err != nil {
						logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot get text content row %d column %d", rowNum, colNum)
						return nil, err
					}
					strAmount = strings.ReplaceAll(strAmount, ",", "")
					strAmount = strings.ReplaceAll(strAmount, "$", "")
					CleanText(&strAmount)
					at.Amount, err = ParseScript(strAmount)
					if err != nil {
						//logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot parse amount: %s, row %d column %d. trying to remove unicode spaces", strAmount, rowNum, colNum)
						CleanTextWithRemoveUnicodeSpaces(&strAmount)
						at.Amount, err = ParseScript(strAmount)
						if err != nil {
							logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot parse amount %s: row %d column %d", strAmount, rowNum, colNum)
							return nil, err
						}
					}
				}

				if colNum == 3 {
					var strTotal string
					strTotal, err = cols[colNum].TextContent()
					if err != nil {
						logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot get text content row %d column %d", r, colNum)
						return nil, err
					}
					strTotal = strings.ReplaceAll(strTotal, ",", "")
					strTotal = strings.ReplaceAll(strTotal, "$", "")
					CleanText(&strTotal)
					at.Total, err = ParseScript(strTotal)
					if err != nil {
						//logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot parse total: %s, row %d column %d. trying to remove unicode spaces", strTotal, rowNum, colNum)
						CleanTextWithRemoveUnicodeSpaces(&strTotal)
						at.Total, err = ParseScript(strTotal)
						if err != nil {
							logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot parse total %s: row %d column %d", strTotal, rowNum, colNum)
							return nil, err
						}
					}
				}

				if colNum == 4 {
					continue
				}

				if colNum == 5 {

					var colSpan = cols[colNum].Locator("xpath=/span")

					ageTimestamp, err := colSpan.GetAttribute("aria-label")
					if err != nil {
						logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot get attribute aria-label: row %d column %d", r, colNum)
						return nil, err
					}

					at.AgeTimestamp, err = time.Parse(time.DateTime, ageTimestamp)
					if err != nil {
						logrus.WithField("wallet_token", walletToken).WithError(err).Errorf("cannot parse age timestamp: row %d column %d", r, colNum)
						return nil, err
					}

				}

				if colNum == 7 {

					var colA = cols[colNum].Locator("xpath=/span/a")
					thxLink, err := colA.GetAttribute("href")
					if err != nil {
						return nil, err
					}

					tlp := strings.Split(thxLink, "/")
					at.WalletTransactionId = tlp[len(tlp)-1]
					//xpathActionsContainer := "//div[@id='wrapperContent']"
					//xpathActions := "xpath=/div/div"
					////np = browser.Contexts()[0].Pages()[0]
					//
					//ac := np.Locator(xpathActionsContainer)
					//allActions, err := ac.Locator(xpathActions).All()
					//if err != nil {
					//	return
					//}
					//
					//var steps [][]string
					//for _, aa := range allActions {
					//	var step []string
					//	txts, err := aa.InnerText()
					//	if err != nil {
					//		return
					//	}
					//	fromString := strings.Split(txts, "\n")[0]
					//
					//	step = append(step, txts)
					//	steps = append(steps, step)
					//}

				}

			}
			ats = append(ats, at)
		}
	}

	return ats, err
}

func locatorWithRetry(page *playwright.Page, locatorPath string, retryCount int) playwright.Locator {
	var locator playwright.Locator
	for i := 0; i < retryCount; i++ {
		err := (*page).Locator(locatorPath).WaitFor(playwright.LocatorWaitForOptions{
			Timeout: playwright.Float(2000),
		})
		if err != nil {
			(*page).WaitForTimeout(1000)
			continue
		} else {
			locator = (*page).Locator(locatorPath)
			break
		}
	}
	return locator
}

func locatorWithRetryCombined(page *playwright.Page, locatorPath string, retryCount int) (playwright.Locator, error) {
	var locator playwright.Locator
	var err error
	for i := 0; i < retryCount; i++ {
		//err := (*page).Locator(locatorPath).WaitFor(playwright.LocatorWaitForOptions{
		//	Timeout: playwright.Float(2000),
		//})
		err = doWaitLocator(page, locatorPath)
		if err != nil {
			(*page).WaitForTimeout(1000)
			continue
		} else {
			locator = (*page).Locator(locatorPath)
			break
		}
	}
	return locator, err
}

func doWaitLocator(page *playwright.Page, locatorPath string) (err error) {
	err = (*page).Locator(locatorPath).WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateAttached,
		Timeout: playwright.Float(2000),
	})
	err = (*page).Locator(locatorPath).WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(2000),
	})

	if err != nil {
		return err
	}

	return nil
}

func doWait(page *playwright.Page, load bool, dom bool, net bool) (err error) {
	if dom {
		err = (*page).WaitForLoadState(playwright.PageWaitForLoadStateOptions{
			State:   playwright.LoadStateDomcontentloaded,
			Timeout: playwright.Float(10000),
		})

	}
	if net {
		err = (*page).WaitForLoadState(playwright.PageWaitForLoadStateOptions{
			State:   playwright.LoadStateNetworkidle,
			Timeout: playwright.Float(10000),
		})
	}
	if load {
		err = (*page).WaitForLoadState(playwright.PageWaitForLoadStateOptions{
			State:   playwright.LoadStateLoad,
			Timeout: playwright.Float(10000),
		})
	}

	if err != nil {
		logrus.WithError(err).Errorf("cannot wait for page")
		return err
	}
	return nil
}
