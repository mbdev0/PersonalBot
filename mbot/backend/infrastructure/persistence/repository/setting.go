package repository

import (
	"context"
	"database/sql"
	"fmt"
	"personal_bot/backend/infrastructure/persistence/mapper"
	"personal_bot/backend/infrastructure/persistence/models"
	"personal_bot/backend/internal/core/settings"
)

type SettingRepository struct {
	db *sql.DB
}

func NewSettingRepository(db *sql.DB) *SettingRepository {
	return &SettingRepository{
		db: db,
	}
}

func (sr *SettingRepository) GetSettings(ctx context.Context) (settings.Settings, error) {
	query := `
	SELECT s.discord_webhook, s.send_on_fail, s.send_on_success, s.position_nodes, s.quick_sell_buttons from settings s`

	row := sr.db.QueryRowContext(ctx, query)

	settingsRow := &models.SettingsRow{}

	err := row.Scan(&settingsRow.DiscordWebhook, &settingsRow.SendOnFail, &settingsRow.SendOnSuccess, &settingsRow.PositionNodes, &settingsRow.QuickSellButtons)
	if err != nil {
		return settings.Settings{}, err
	}

	coreSettings, err := mapper.MapRepoToSettings(settingsRow)
	if err != nil {
		return settings.Settings{}, err
	}

	return coreSettings, nil

}

func (sr *SettingRepository) PostSettings(ctx context.Context, settings settings.Settings) error {
	query := `
		UPDATE settings
		SET discord_webhook=?, send_on_fail=?, send_on_success=?, position_nodes=?, quick_sell_buttons=?
		WHERE id = 1
	`
	mappedSettings, err := mapper.MapSettingsToRepo(settings)
	if err != nil {
		return err
	}

	row, err := sr.db.ExecContext(ctx, query,
		mappedSettings.DiscordWebhook, mappedSettings.SendOnFail, mappedSettings.SendOnSuccess,
		mappedSettings.PositionNodes, mappedSettings.QuickSellButtons)

	if err != nil {
		return err
	}

	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 1 {
		return nil
	} else {
		return fmt.Errorf("no rows affected - make sure ")
	}
}
