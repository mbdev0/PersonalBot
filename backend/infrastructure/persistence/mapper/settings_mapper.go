package mapper

import (
	"encoding/json"
	"fmt"
	"personal_bot/infrastructure/persistence/models"
	"personal_bot/internal/core/settings"
)

func MapSettingsToRepo(settings settings.Settings) (settingsRow models.SettingsRow, err error) {
	settingsRow.DiscordWebhook = settings.DiscordWebhook
	if settings.SendOnFail {
		settingsRow.SendOnFail = 1
	} else {
		settingsRow.SendOnFail = 0
	}

	if settings.SendOnSuccess {
		settingsRow.SendOnSuccess = 1
	} else {
		settingsRow.SendOnSuccess = 0
	}

	posNodes, err := json.Marshal(settings.PositionNodes)
	if err != nil {
		return
	}
	settingsRow.PositionNodes = string(posNodes)

	quickSellButtons, err := json.Marshal(settings.QuickSellButtons)
	if err != nil {
		return
	}

	settingsRow.QuickSellButtons = string(quickSellButtons)

	return settingsRow, nil
}

func MapRepoToSettings(settingsRow *models.SettingsRow) (sttings settings.Settings, err error) {
	sttings.DiscordWebhook = settingsRow.DiscordWebhook
	sttings.SendOnFail, err = getBoolean(settingsRow.SendOnFail)
	if err != nil {
		return sttings, err
	}
	sttings.SendOnSuccess, err = getBoolean(settingsRow.SendOnSuccess)
	if err != nil {
		return sttings, err
	}

	positionNodes := new(settings.PositionNodes)
	err = json.Unmarshal([]byte(settingsRow.PositionNodes), positionNodes)
	if err != nil {
		return
	}
	sttings.PositionNodes = *positionNodes

	quickSellButtons := new(settings.QuickSellButtons)
	err = json.Unmarshal([]byte(settingsRow.QuickSellButtons), quickSellButtons)
	if err != nil {
		return
	}

	sttings.QuickSellButtons = *quickSellButtons

	return sttings, nil
}

func getBoolean(num int64) (bool, error) {
	switch num {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, fmt.Errorf("number is not boolean - check settings repo")
	}
}
