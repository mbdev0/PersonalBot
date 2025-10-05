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
		logger.Error("Error reading userConfig file", err)
		return
	}

	userConfig := config.GetConfig()
	fmt.Println("HTTP Node: ", userConfig.HttpNode)
	fmt.Println("WS Node: ", userConfig.WsNode)
	fmt.Println("Webhook: ", userConfig.Webhook)

	pumpfun_idl.GetIdlMap()
	_, err = solana_price.GetSolPrice()
	if err != nil {
		return
	}
	lookuptable.GetAddressLookupTable()
	validator.GetValidator()
}
