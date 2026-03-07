package mapper

import (
	"fmt"
	"personal_bot/infrastructure/persistence/models"
	"personal_bot/internal/core/settings"
)

func MapSettingsToRepo(settings settings.Settings) (settingsRow models.SettingsRow) {
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

	return settingsRow
}

func MapRepoToSettings(settingsRow *models.SettingsRow) (settings settings.Settings, err error) {
	settings.DiscordWebhook = settingsRow.DiscordWebhook
	settings.SendOnFail, err = getBoolean(settingsRow.SendOnFail)
	if err != nil {
		return settings, err
	}
	settings.SendOnSuccess, err = getBoolean(settingsRow.SendOnSuccess)
	if err != nil {
		return settings, err
	}

	return settings, nil
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
