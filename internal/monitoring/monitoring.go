package monitoring

import (
	"pump_fun/internal/handlers"
	"pump_fun/internal/logger"
	"pump_fun/internal/models"
	"pump_fun/internal/monitoring/geyser"
	"pump_fun/internal/monitoring/transactions/decoder"
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
			transaction_notification_chan := make(chan models.TransactionNotification, 1000)
			coinStructChan := make(chan models.Coin, 1000)

			d := &decoder.Decoder{}
			createInstructionDecoder := &decoder.CreateInstructionDecoder{}
			d.SetDecodingStrategy(createInstructionDecoder)

			go func() {
				err := geyser.Geyser_Stream_Transactions(transaction_notification_chan)
				if err != nil {
					logger.Error(err)
				}
			}()

			go func() {
				for transaction := range transaction_notification_chan {
					go func() {
						handlers.HandleTransactionNotification(d, transaction, coinStructChan)
					}()
				}

			}()

			for coinStruct := range coinStructChan {
				logger.Log(logger.LevelInfo, "Coin Struct: ", logger.String("Coin", coinStruct.CoinData.Name))
				go func() {

					webhook.SendWebhook(&coinStruct)
				}()
			}

			wg.Done()
		}()
	}

	wg.Wait()
}
