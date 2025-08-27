package config

import (
	"encoding/json"
	"os"
	"pump_fun/internal/core/models"
	"sync"
)

var (
	config *models.Config
	once   sync.Once
)

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
	once.Do(func() {
		err := LoadConfig()
		if err != nil {
			return
		}
	})
	return config
}
