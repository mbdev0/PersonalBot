package monitoring

import (
	"pump_fun/internal/monitoring/decoder"
	"pump_fun/internal/monitoring/filters"
	"pump_fun/internal/monitoring/models"
	"pump_fun/internal/monitoring/stream"
	"pump_fun/internal/monitoring/stream/response"
	"pump_fun/pkg/logger"
	"sync"
)

// temporary bool until gui is built
var startMonitoring = true
var transactionChanSize = 1000
var coinChanSize = 1000

func StartAFKMonitor(coins chan<- models.Coin) {
	var wg sync.WaitGroup
	if startMonitoring {
		wg.Add(1)
		go func() {
			defer wg.Done()
			transactionNotificationChan := make(chan response.TransactionNotification, transactionChanSize)
			coinStructChan := make(chan models.Coin, coinChanSize)
			defer close(transactionNotificationChan)
			defer close(coinStructChan)

			go streamTransactions(transactionNotificationChan)

			go decryptAndFilterTransactions(transactionNotificationChan, coinStructChan)

			for coinStruct := range coinStructChan {
				coins <- coinStruct
				logger.Information("coin found: " + coinStruct.CoinData.Symbol)
				// go func(coin models.Coin) {
				// 	logger.Information("Coin found: " + coin.CoinData.Symbol)
				// 	// webhook.SendWebhook(&coin)
				// }(coinStruct)
			}
		}()
	}

	wg.Wait()
}

// was splitting this into functions a good idea?
// -> harder to follow the "pipeline of what's going on and what being passed"
// -> cleans up main function
// TODO: review the above
func streamTransactions(transactionNotificationChan chan<- response.TransactionNotification) {
	err := stream.GeyserStreamTransactions(transactionNotificationChan)
	if err != nil {
		logger.Error("Error in Geyser_Stream_Transactions ", err)
		close(transactionNotificationChan)
	}
}

func decryptAndFilterTransactions(transactionNotificationChan <-chan response.TransactionNotification, coinStructChan chan<- models.Coin) {
	defer close(coinStructChan)
	for transaction := range transactionNotificationChan {
		go func(transaction response.TransactionNotification, coinStructChan chan<- models.Coin) {
			handleTransactionNotification(transaction, coinStructChan)
		}(transaction, coinStructChan)
	}
}

func handleTransactionNotification(transaction response.TransactionNotification, coinStructChan chan<- models.Coin) {
	coin := decoder.DecryptTransactionNotificationForCoin(transaction)
	if coin != nil {
		coin = filters.HandleCoinFiltering(coin) //we should ideally get rid of this from here and do a filter.ApplyFilters() here
	}

	if coin != nil {
		coinStructChan <- *coin
	}
}
