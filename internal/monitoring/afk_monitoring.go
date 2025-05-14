package monitoring

import (
	"pump_fun/internal/handlers"
	"pump_fun/internal/logger"
	"pump_fun/internal/models"
	"pump_fun/internal/monitoring/geyser"
	"pump_fun/internal/webhook"
	"sync"
)

// temporary bool until gui is built
var startMonitoring bool = true

func StartAFKMonitor() {
	var wg sync.WaitGroup

	if startMonitoring {
		wg.Add(1)
		go func() {
			transaction_notification_chan := make(chan models.TransactionNotification, 1000)
			coinStructChan := make(chan models.Coin, 1000)
      
			go MonitorTransactions(transaction_notification_chan)
			go ProcessAndFilterTransactions(transaction_notification_chan, coinStructChan)

      for coinStruct := range coinStructChan {
				go func(coin models.Coin) {
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
		logger.Log(logger.LevelError, "Error in Geyser_Stream_Transactions", logger.Error(err))
	}
}

func ProcessAndFilterTransactions(transaction_notification_chan <-chan models.TransactionNotification, coinStructChan chan<- models.Coin) {
	for transaction := range transaction_notification_chan {
		go func(transaction models.TransactionNotification, coinStructChan chan<- models.Coin) {
			handlers.HandleTransactionNotification(transaction, coinStructChan)
		}(transaction, coinStructChan)
	}
}
