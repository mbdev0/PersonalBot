package handlers

import (
	"pump_fun/internal/models"
	"pump_fun/internal/monitoring/transactions/transaction_decoder"
)

func HandleTransactionNotification(transaction models.TransactionNotification, coinStructChan chan models.Coin) {
	coin := transaction_decoder.DecryptTransactionNotificationForCoin(transaction, coinStructChan)

	if coin != nil {
		coin = HandleCoinFiltering(coin)
	}

	if coin != nil {
		coinStructChan <- *coin
	}
}
