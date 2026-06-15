package models

type Config struct {
	Webhook          string `json:"webhook"`
	HttpNode         string `json:"http_node"`
	WsNode           string `json:"ws_node"`
	WalletPrivateKey string `json:"wallet_privatekey"`
}