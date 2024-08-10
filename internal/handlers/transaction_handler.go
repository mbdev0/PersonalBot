package handlers

import (
	"pump_fun/internal/monitoring/transactions"
	"pump_fun/internal/webhook"
)

func HandleTransaction(signature string) {
	transactions.GetTransaction(signature)
	webhook.SendTelegramMessage("", "", "")
}
