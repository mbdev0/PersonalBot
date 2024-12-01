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

func StartMonitor() {
	var wg sync.WaitGroup

	// start monitoring on 1 goroutine
	if startMonitoring {
		wg.Add(1)
		go func() {
			transaction_notification_chan := make(chan geyser.TransactionNotification, 1000)
			coinStructChan := make(chan models.Coin, 1000)

			go func() {
				err := geyser.Geyser_Stream_Transactions(transaction_notification_chan)
				if err != nil {
					logger.Error(err)
				}
			}()

			// we should then pass the transaction to a handler which will decrypt the transaction and return a struct

			go func() {
				for transaction := range transaction_notification_chan {
					go func() {
						handlers.HandleTransactionNotification(transaction, coinStructChan)
					}()
				}

			}()

			for coinStruct := range coinStructChan {
				go func() {

					webhook.SendWebhook(&coinStruct)
				}()
			}

			wg.Done()
		}()
	}

	wg.Wait()
}
