package sync

import (
	db_wt "order-executor/pkg/model/database/wallet_transaction"
	"time"

	"gorm.io/gorm"
)

type Env struct {
	MinPercOrderThreshold        int
	MaxPercOrderThreshold        int
	SetMaxPercThreshold          bool
	SetMinPercThreshold          bool
	OrderTimeExpirationThreshold int
	MaxPriceRangePerc            float32
	Database                     *gorm.DB
}

func (e *Env) ExecuteOrdres() (response interface{}, err error) {

	// prendi tutte le transazioni recuperate non ancora processate

	var cryptoTransactions []db_wt.WalletTransaction
	err = e.Database.Where("ProcessedByOrderExecutor = ?", false).Preload("Wallet.WalletToken").Find(&cryptoTransactions).Error

	if err != nil {
		return nil, err
	}

	if len(cryptoTransactions) > 0 {
		currentProcessDate := time.Now()
		expiredThresholdParameter := time.Duration(e.OrderTimeExpirationThreshold) * time.Minute
		expiredThresholdDate := currentProcessDate.Add(-expiredThresholdParameter)
		for _, ct := range cryptoTransactions {

			// 	// verifica che la transazione non sia scaduta
			if ct.CreatedAt.Before(expiredThresholdDate) {
				// scaduta
				// marcala solo come processata
				// forse va aggiunto un campo di stato di processamento scaduta, aperta, chiusa, annullata per superamento soglia ecc
			}
			// prendo il wallet e poi il wallet token amount () usa TokenId di transactio.token per combaciare con wallet_token.
			ct.Wallet
			// 	// prendi il saldo del wallet_token da cui parte la transazione,
			// 	// calcola la percentuale

			// il mmio saldo lo prendo mio wallet ()  c'è il type, se type = 0 allora è il mio. type 1 sono gli altri.
			// per i compra, posso non avere quel token nel mio wallet.
			// quindi devo comprare da eth verso il token di destinazione
			// devo capirese lui va da A --> B, io devo andare da ETH a B, quindi devo capire quanti ETH vale il suo cambio.  ho l'ammontare in dollari, allora calcolo da dollari ad eth. mi serve un'api
			//

			// se vedo una transazione in vendità, verifico se ho il token che lui sta vendendo, se non ce lo ho allora nulla,
			// se ce lo ho allora lo vendo

			// quindi lui compra X allora io compro da eth ad X.
			// quando lui vende, se ho quel token, lo vendo. se non ce lo ho passo
			// quando invece vendo io perchè mi basta, allora mi serve il valore in dollari corrente del token.

			// limitare il numero di transazioni attive e vendere quando hai un tot di guadagno, ma se lui vende prima del guadagno allora vendi
			// stessa cosa per la perdita.
			// 	// calcola il proprio saldo sullo stesso token.
			// se non ho fondi in quel token salto
			// 	// calcola la percentuale riportata al proprio saldo

			// 	// verifica che rispetti lee regole di threshold
		}
	}

	// genera il corrispettivo ordine

	// salva la transazione come processata

	return nil, nil
}
