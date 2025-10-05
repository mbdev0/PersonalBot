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

			var handlerWg sync.WaitGroup

			go streamTransactions(transactionNotificationChan, ctx)

			go decryptAndFilterTransactions(filters, transactionNotificationChan, coinStructChan, ctx, &handlerWg)

			for {
				select {
				case <-ctx.Done():
					handlerWg.Wait()
					close(transactionNotificationChan)
					close(coinStructChan)
					return
				case coin, ok := <-coinStructChan:
					if !ok {
						return
					}
					coins <- coin
				}
			}

		}()
	}

	wg.Wait()
}

func streamTransactions(transactionNotificationChan chan<- response.TransactionNotification, ctx context.Context) {
	err := stream.StartGeyserTransactionStream(transactionNotificationChan, ctx)
	if err != nil {
		logger.Error("Error in Geyser_Stream_Transactions ", err)
		return
	}
}

func decryptAndFilterTransactions(filters filters.FilterPipeline, transactionNotificationChan <-chan response.TransactionNotification, coinStructChan chan<- models.Coin, ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			return
		case transaction, ok := <-transactionNotificationChan:
			if !ok {
				return
			}
			wg.Add(1)
			go func(transaction response.TransactionNotification) {
				defer wg.Done()
				handleTransactionNotification(filters, transaction, coinStructChan, ctx)
			}(transaction)
		}
	}
}

func handleTransactionNotification(filters filters.FilterPipeline, transaction response.TransactionNotification, coinStructChan chan<- models.Coin, ctx context.Context) {
	coin := decoder.DecryptTransactionNotificationForCoin(transaction)
	if coin != nil {
		coin = filters.ApplyFilters(coin)
	}

	if coin != nil {
		coinStructChan <- *coin
	}
}
