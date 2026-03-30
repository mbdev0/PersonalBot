package settings

type Settings struct {
	DiscordWebhook string        `json:"discord_webhook"`
	SendOnFail     bool          `json:"send_on_fail"`
	SendOnSuccess  bool          `json:"send_on_success"`
	PositionNodes  PositionNodes `json:"position_nodes"`
}

type PositionNodes struct {
	HTTPNode string `json:"http_node"`
	WSNode   string `json:"ws_node"`
}

type TestWebhook struct {
	DiscordWebhook string `json:"discord_webhook"`
}
