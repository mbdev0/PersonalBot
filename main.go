// temporary file for go build
package main

import (
	"fmt"
	"pump_fun/internal/monitoring"
)

func main() {
	fmt.Println("Hello, World!")

	monitoring.StartMonitor()
}
