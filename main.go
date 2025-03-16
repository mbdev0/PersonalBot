// temporary file for go build
package main

import (
	"pump_fun/internal/launch"
	"pump_fun/internal/monitoring"
)

func main() {
	launch.LaunchOperations()
	monitoring.StartAFKMonitor()
}
