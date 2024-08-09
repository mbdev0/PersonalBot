package monitoring

import (
	"pump_fun/internal/monitoring/logsubscribe"
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
			logsubscribe.LogSubscribe()
			wg.Done()
		}()
	}

	wg.Wait()
}
