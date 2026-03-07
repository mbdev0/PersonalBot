package controller

import (
	"context"
	models "personal_bot/internal/core/settings"
	"personal_bot/internal/services/settings"
)

type SettingsController struct {
	service *settings.Service
}

func NewSettingsController(service *settings.Service) *SettingsController {
	return &SettingsController{
		service: service,
	}
}

func (sc *SettingsController) GetSettings(ctx context.Context) (models.Settings, error) {
	return sc.service.GetSettings(ctx)
}

func (sc *SettingsController) PostSettings(ctx context.Context, settings models.Settings) error {
	return sc.service.PostSettings(ctx, settings)
}

func (sc *SettingsController) TestDiscordEndpoint() {
}
