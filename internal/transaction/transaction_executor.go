package transaction

import (
	"pump_fun/internal/models/tasks"
	"pump_fun/internal/transaction/buy"
	"pump_fun/internal/transaction/sell"
)

type TransactionExecutor struct{}

func (th *TransactionExecutor) Execute(task tasks.Task) {
	switch t := task.(type) {
	case *tasks.BuyTask:
		go buy.SendBuyTransaction(t)
	case *tasks.SellTask:
		go sell.SendSellTransaction(t)
	}
}
