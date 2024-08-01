package handlers

// TODO : Check these import did not test
import (
    "pump_fun/internal/transactions/"
    "pump_fun/internal/webhook"
)

func HandleTransaction(signature string) {
    transactionData := get_transaction.GetTransaction(signature)
    webhook.SendToWebhook(transactionData)
}