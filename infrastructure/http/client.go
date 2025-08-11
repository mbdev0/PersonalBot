package http

import (
	"net/http"
	"sync"
	"time"
)

var lock = &sync.Mutex{}
var connectionTimeout = time.Second * 10
var client *http.Client

func GetClient() *http.Client {
	if client == nil {
		lock.Lock()
		defer lock.Unlock()

		if client == nil {
			client = &http.Client{Timeout: connectionTimeout}
		}
	}

	return client
}
