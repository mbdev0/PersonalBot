package monitoring

import (
	"context"
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

func StartAFKMonitor(filters filters.FilterPipeline, coins chan<- models.Coin, ctx context.Context) {
	var wg sync.WaitGroup
	if startMonitoring {
		wg.Add(1)
		go func() {
			defer wg.Done()
			transactionNotificationChan := make(chan response.TransactionNotification, transactionChanSize)
			coinStructChan := make(chan models.Coin, coinChanSize)
			defer close(transactionNotificationChan)
			defer close(coinStructChan)

			go streamTransactions(transactionNotificationChan, ctx)

			go decryptAndFilterTransactions(filters, transactionNotificationChan, coinStructChan, ctx)

			for coinStruct := range coinStructChan {
				coins <- coinStruct
				logger.Information("coin found: " + coinStruct.CoinData.Symbol)
			}
		}()
	}

	wg.Wait()
}

func streamTransactions(transactionNotificationChan chan<- response.TransactionNotification, ctx context.Context) {
	err := stream.StartGeyserTransactionStream(transactionNotificationChan, ctx) //pass in ctx here
	if err != nil {
		logger.Error("Error in Geyser_Stream_Transactions ", err)
		return
	}
}

func decryptAndFilterTransactions(filters filters.FilterPipeline, transactionNotificationChan <-chan response.TransactionNotification, coinStructChan chan<- models.Coin, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case transaction := <-transactionNotificationChan:
			go func(transaction response.TransactionNotification) {
				handleTransactionNotification(filters, transaction, coinStructChan)
			}(transaction)
		}
	}
}

func handleTransactionNotification(filters filters.FilterPipeline, transaction response.TransactionNotification, coinStructChan chan<- models.Coin) {
	coin := decoder.DecryptTransactionNotificationForCoin(transaction)
	if coin != nil {
		coin = filters.ApplyFilters(coin)
	}

	if coin != nil {
		coinStructChan <- *coin
	}
}
