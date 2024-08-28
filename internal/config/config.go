package config

import (
	"encoding/json"
	"os"
	"pump_fun/internal/logger"
)

var config *Config

func LoadConfig() error {
	file, err := os.ReadFile("configuration/config.json")
	if err != nil {
		logger.Log(logger.LevelError, "Error reading config file", logger.String("error: ", err.Error()))
		return err
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		logger.Log(logger.LevelError, "Error unmarshalling config file", logger.String("error: ", err.Error()))
		return err
	}
	return nil
}

func GetConfig() *Config {
	return config
}
