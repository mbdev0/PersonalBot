// temporary file for go build
package main

import (
	"fmt"
	"pump_fun/internal/monitoring/logsubscribe"
	"sync"
)

func main() {
	fmt.Println("Hello, World!")
	var wg sync.WaitGroup
	userClickedStart := true

	// run the log subscribe function in a goroutine and wait for it to finish
	// the idea is that this function will run when a user clicks a button on the frontend
	if userClickedStart {
		wg.Add(1)
		go func() {
			logsubscribe.LogSubscribe()
			wg.Done()
		}()
	}

	wg.Wait()

}
