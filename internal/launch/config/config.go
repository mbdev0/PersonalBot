package config

import (
	"encoding/json"
	"os"
	"pump_fun/internal/logger"
	"pump_fun/internal/models"
	"sync"
)

var config *models.Config
var lock = &sync.Mutex{}

func LoadConfig() error {
	file, err := os.ReadFile("configuration/config.json")
	if err != nil {
		logger.Log(logger.LevelError, "Error reading config file", logger.Error(err))
		return err
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		logger.Log(logger.LevelError, "Error unmarshalling config file", logger.Error(err))
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
				logger.Log(logger.LevelError, "Error loading config", logger.Error(err))
				return &models.Config{}
			}
		}
	}

	return config
}
