package app

import (
	"fmt"
	lookuptable "pump_fun/app/lookup_table"
	"pump_fun/app/pumpfun_idl"
	"pump_fun/app/validator"
	"pump_fun/infrastructure/config"
	"pump_fun/infrastructure/solana_price"
	"pump_fun/pkg/logger"
)

func Launch() {
	err := config.LoadConfig()
	if err != nil {
		logger.Error("Error reading config file", err)
		return
	}

	config := config.GetConfig()
	fmt.Println("HTTP Node: ", config.HttpNode)
	fmt.Println("WS Node: ", config.WsNode)
	fmt.Println("Webhook: ", config.Webhook)

	pumpfun_idl.GetIdlMap()
	solana_price.GetSolPrice()
	lookuptable.GetAddressLookupTable()
	validator.GetValidator()
}
