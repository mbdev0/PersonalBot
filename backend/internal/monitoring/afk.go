package monitoring

import (
	"context"
	"personal_bot/internal/monitoring/decoder"
	"personal_bot/internal/monitoring/filters"
	"personal_bot/internal/monitoring/models"
	"personal_bot/internal/monitoring/stream"
	"personal_bot/internal/monitoring/stream/response"
	"personal_bot/pkg/logger"
	"sync"
)

var startMonitoring = true
var transactionChanSize = 1000
var coinChanSize = 1000

func StartAFKMonitor(ctx context.Context, filters filters.FilterPipeline, coins chan<- models.Coin, wsUrl string) {
	var wg sync.WaitGroup
	if startMonitoring {
		wg.Add(1)
		go func() {
			defer wg.Done()
			transactionNotificationChan := make(chan response.TransactionNotification, transactionChanSize)
			coinStructChan := make(chan models.Coin, coinChanSize)

			var handlerWg sync.WaitGroup

			go streamTransactions(ctx, transactionNotificationChan, wsUrl)

			go decryptAndFilterTransactions(ctx, filters, transactionNotificationChan, coinStructChan, &handlerWg)

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

func streamTransactions(ctx context.Context, transactionNotificationChan chan<- response.TransactionNotification, wsUrl string) {
	err := stream.StartGeyserTransactionStream(ctx, transactionNotificationChan, wsUrl)
	if err != nil {
		logger.Error("Error in Geyser_Stream_Transactions ", err)
		return
	}
}

func decryptAndFilterTransactions(ctx context.Context, filters filters.FilterPipeline, transactionNotificationChan <-chan response.TransactionNotification, coinStructChan chan<- models.Coin, wg *sync.WaitGroup) {
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
