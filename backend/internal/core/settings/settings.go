package settings

type Settings struct {
	DiscordWebhook   string           `json:"discord_webhook"`
	SendOnFail       bool             `json:"send_on_fail"`
	SendOnSuccess    bool             `json:"send_on_success"`
	PositionNodes    PositionNodes    `json:"position_nodes"`
	QuickSellButtons QuickSellButtons `json:"quick_sell_buttons"`
}

type PositionNodes struct {
	HTTPNode string `json:"http_node"`
	WSNode   string `json:"ws_node"`
}

type QuickSellButtons struct {
	Button1 float64 `json:"button_1"`
	Button2 float64 `json:"button_2"`
	Button3 float64 `json:"button_3"`
	Button4 float64 `json:"button_4"`
}

type TestWebhook struct {
	DiscordWebhook string `json:"discord_webhook"`
}
