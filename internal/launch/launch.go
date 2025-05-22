package launch

import (
	"fmt"
	"pump_fun/internal/launch/config"
	"pump_fun/internal/launch/pumpfun_idl"
	"pump_fun/internal/launch/solana_price"
	"pump_fun/internal/logger"
)

func LaunchOperations() {
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
}
