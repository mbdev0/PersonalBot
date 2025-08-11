package config

import (
	"encoding/json"
	"os"
	"pump_fun/internal/core/models"
	"pump_fun/pkg/logger"
	"sync"
)

var config *models.Config
var lock = &sync.Mutex{}

func LoadConfig() error {
	file, err := os.ReadFile("configuration/config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig() *models.Config {
	if config == nil {
		lock.Lock()
		defer lock.Unlock()
		if config == nil {
			err := LoadConfig()
			if err != nil {
				logger.Error("Error loading config", err)
				return &models.Config{}
			}
		}
	}

	return config
}
