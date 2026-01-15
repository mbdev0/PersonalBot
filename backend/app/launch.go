package app

import (
	"fmt"
	lookuptable "personal_bot/app/lookup_table"
	"personal_bot/app/pumpfun_idl"
	"personal_bot/app/validator"
	"personal_bot/infrastructure/config"
	"personal_bot/infrastructure/solana_price"
	"personal_bot/pkg/logger"
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
