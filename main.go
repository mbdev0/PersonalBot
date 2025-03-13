// temporary file for go build
package main

import (
	"context"
	"fmt"
	"pump_fun/internal/config"
	"pump_fun/internal/launch"
	"pump_fun/internal/monitoring"
	"time"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		return
	}

	config := config.GetConfig()
	fmt.Println("HTTP Node: ", config.HttpNode)
	fmt.Println("WS Node: ", config.WsNode)
	fmt.Println("Webhook: ", config.Webhook)

	launch.GetSolPrice()
	launch.GetIdlMap()

	// monitoring.StartAFKMonitor()
	ctx, cancel := context.WithCancel(context.Background())
	address := "5U3smn2USzQGSWfM3JmKXt9YPKpWifjxvR2aMZkQAN1S"
	go monitoring.StartMarketCapMonitor(ctx, address)

	time.Sleep(10 * time.Second)
	cancel()
}
