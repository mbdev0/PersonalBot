// temporary file for go build
package main

import (
	"fmt"
	"pump_fun/internal/config"
	"pump_fun/internal/handlers"
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

	// monitoring.StartMonitor()

	handlers.HandleTransaction("3CahbK6dzgHFzzCKLe8yk9WTeYPUmBRDte9TzQNef489yoYbaxkxtjnpg5Y2EJ9au99gUu1zozZ6ysMYQq9sDUSG")
}
