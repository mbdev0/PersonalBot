package launch

import (
	"fmt"
	"pump_fun/internal/launch/config"
	"pump_fun/internal/launch/pumpfun_idl"
	"pump_fun/internal/launch/solana_price"
)

func LaunchOperations() {
	err := config.LoadConfig()
	if err != nil {
		return
	}

	config := config.GetConfig()
	fmt.Println("HTTP Node: ", config.HttpNode)
	fmt.Println("WS Node: ", config.WsNode)
	fmt.Println("Webhook: ", config.Webhook)

	pumpfun_idl.GetIdlMap()
	solana_price.GetSolPrice()
}
