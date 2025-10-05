package http

import (
	"net/http"
	"sync"
	"time"
)

var (
	once              sync.Once
	connectionTimeout = time.Second * 10
	client            *http.Client
)

func GetClient() *http.Client {
	once.Do(func() {
		client = &http.Client{Timeout: connectionTimeout}
	})
	return client
}
