package controller

import (
	"context"
	models "personal_bot/backend/internal/core/settings"
	"personal_bot/backend/internal/services/notifier"
	"personal_bot/backend/internal/services/settings"
)

type SettingsController struct {
	service  *settings.Service
	notifier *notifier.DiscordNotifier
}

func NewSettingsController(service *settings.Service, notifier *notifier.DiscordNotifier) *SettingsController {
	return &SettingsController{
		service:  service,
		notifier: notifier,
	}
}

func (sc *SettingsController) GetSettings(ctx context.Context) (models.Settings, error) {
	return sc.service.GetSettings(ctx)
}

func (sc *SettingsController) PostSettings(ctx context.Context, settings models.Settings) error {
	sc.notifier.Update(settings)
	return sc.service.PostSettings(ctx, settings)
}

func (sc *SettingsController) TestDiscordEndpoint(discordWebhookUrl string) error {
	return sc.notifier.TestDiscordEndpoint(discordWebhookUrl)
}
