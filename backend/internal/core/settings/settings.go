package settings

type Settings struct {
	DiscordWebhook string `json:"discord_webhook"`
	SendOnFail     bool   `json:"send_on_fail"`
	SendOnSuccess  bool   `json:"send_on_success"`
}

type TestWebhook struct {
	DiscordWebhook string `json:"discord_webhook"`
}
