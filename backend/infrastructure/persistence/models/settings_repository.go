package models

type SettingsRow struct {
	DiscordWebhook string `db:"discord_webhook"`
	SendOnFail     int64  `db:"send_on_fail"`
	SendOnSuccess  int64  `db:"send_on_success"`
	PositionNodes  string `db:"position_nodes"`
}
