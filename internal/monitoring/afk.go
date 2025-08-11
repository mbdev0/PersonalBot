package monitoring

import (
	"pump_fun/infrastructure/webhook"
	"pump_fun/internal/monitoring/decoder"
	"pump_fun/internal/monitoring/filters"
	"pump_fun/internal/monitoring/models"
	"pump_fun/internal/monitoring/stream"
	"pump_fun/internal/monitoring/stream/response"
	"pump_fun/pkg/logger"
	"sync"
)

// temporary bool until gui is built
var startMonitoring bool = true
var transaction_chan_size = 1000
var coin_chan_size = 1000

func StartAFKMonitor() {
	var wg sync.WaitGroup
	if startMonitoring {
		wg.Add(1)
		go func() {
			transaction_notification_chan := make(chan response.TransactionNotification, transaction_chan_size)
			coinStructChan := make(chan models.Coin, coin_chan_size)

			go MonitorTransactions(transaction_notification_chan)

			go ProcessAndFilterTransactions(transaction_notification_chan, coinStructChan)

			for coinStruct := range coinStructChan {
				go func(coin models.Coin) {
					logger.Information("Coin found: " + coin.CoinData.Name)
					webhook.SendWebhook(&coin)
				}(coinStruct)
			}

			wg.Done()
		}()
	}

	wg.Wait()
}

func MonitorTransactions(transaction_notification_chan chan<- response.TransactionNotification) {
	err := stream.Geyser_Stream_Transactions(transaction_notification_chan)
	if err != nil {
		logger.Error("Error in Geyser_Stream_Transactions ", err)
		close(transaction_notification_chan)
	}
}

func ProcessAndFilterTransactions(transaction_notification_chan <-chan response.TransactionNotification, coinStructChan chan<- models.Coin) {
	defer close(coinStructChan)
	for transaction := range transaction_notification_chan {
		go func(transaction response.TransactionNotification, coinStructChan chan<- models.Coin) {
			handleTransactionNotification(transaction, coinStructChan)
		}(transaction, coinStructChan)
	}
}

func handleTransactionNotification(transaction response.TransactionNotification, coinStructChan chan<- models.Coin) {
	coin := decoder.DecryptTransactionNotificationForCoin(transaction)

	if coin != nil {
		coin = filters.HandleCoinFiltering(coin)
	}

	if coin != nil {
		coinStructChan <- *coin
	}
}
