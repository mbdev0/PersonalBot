package monitoring

import (
	"fmt"
	"pump_fun/internal/monitoring/geyser"
	"sync"
)

// temporary bool until gui is built
var startMonitoring bool = true

func StartMonitor() {
	var wg sync.WaitGroup

	// start monitoring on 1 goroutine
	if startMonitoring {
		wg.Add(1)
		go func() {
			err := geyser.Geyser_Stream_Transactions()
			// err := logsubscribe.LogSubscribe()
			if err != nil {
				// should be updated to show the error on the ui
				fmt.Println("Error: ", err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
