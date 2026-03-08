package settings

import (
	"context"
	"personal_bot/infrastructure/persistence/repository"
	"personal_bot/internal/core/settings"
)

type Service struct {
	repo *repository.SettingRepository
}

func NewSettingsService(repo *repository.SettingRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetSettings(ctx context.Context) (settings.Settings, error) {
	return s.repo.GetSettings(ctx)
}

func (s *Service) PostSettings(ctx context.Context, settings settings.Settings) error {
	return s.repo.PostSettings(ctx, settings)
}
