package monitoring

import (
	"pump_fun/internal/handlers"
	"pump_fun/internal/models"
	"pump_fun/internal/monitoring/geyser"
	"pump_fun/internal/webhook"
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
			transaction_notification_chan := make(chan models.TransactionNotification, transaction_chan_size)
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

func MonitorTransactions(transaction_notification_chan chan<- models.TransactionNotification) {
	err := geyser.Geyser_Stream_Transactions(transaction_notification_chan)
	if err != nil {
		logger.Error("Error in Geyser_Stream_Transactions ", err)
		close(transaction_notification_chan)
	}
}

func ProcessAndFilterTransactions(transaction_notification_chan <-chan models.TransactionNotification, coinStructChan chan<- models.Coin) {
	defer close(coinStructChan)
	for transaction := range transaction_notification_chan {
		go func(transaction models.TransactionNotification, coinStructChan chan<- models.Coin) {
			handlers.HandleTransactionNotification(transaction, coinStructChan)
		}(transaction, coinStructChan)
	}
}
