package webhook

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"encoding/json"
	"net/http"

	"pump_fun/internal/config"
	"pump_fun/internal/logger"
	"pump_fun/internal/models"
)

var (
	discordWebhookURL = config.GetConfig().Webhook
)

func SendWebhook(coin *models.Coin) {

	err := sendDiscordMessage(discordWebhookURL, *coin)

	if err != nil {
		logger.Log(logger.LevelError, "Error sending discord message", logger.Error(err))
	}
}

func sendDiscordMessage(webhookURL string, coin models.Coin) error {
	transport := &http.Transport{
		IdleConnTimeout: 30 * time.Second,
	}
	client := &http.Client{Transport: transport}

	webhook := formatCoinInfo(coin)
	marshaledWebhook, err := json.Marshal(webhook)

	if err != nil {
		logger.Log(logger.LevelError, "Error marshalling webhook", logger.Error(err))
		return err
	}

	req, err := client.Post(webhookURL, "application/json", bytes.NewBuffer(marshaledWebhook))

	if err != nil {
		logger.Log(logger.LevelError, "Error sending webhook", logger.Error(err))
		return err
	}

	defer req.Body.Close()

	if req.StatusCode != 204 {

		body, err := io.ReadAll(req.Body)

		if err != nil {
			logger.Log(logger.LevelError, "Error reading response body", logger.Error(err))
			return err
		}

		if err := handleError(req, body); err != nil {
			return err
		}
	}

	return nil
}

func formatCoinInfo(coin models.Coin) Webhook {
	embed := Embeds{
		Title:  fmt.Sprintf("%s | %s ", coin.CoinData.Name, coin.CoinData.Symbol),
		URL:    "https://pump.fun/" + coin.CoinData.TokenAddr,
		Color:  5814783,
		Fields: generateFields(coin),
		Author: Author{
			Name: "New Pairs Monitor",
		},
		Thumbnail: Thumbnail{
			URL: coin.IPFSData.ImageURL,
		},
	}
	webhook := Webhook{
		Embeds: []Embeds{embed},
	}
	return webhook

}

func generateFields(coin models.Coin) []Fields {
	socialsValue := buildSocials(coin.IPFSData)

	fields := []Fields{
		{
			Name:  "Mint Address",
			Value: fmt.Sprintf("`%s`", coin.CoinData.TokenAddr),
		},
		{
			Name:  "Creator Address",
			Value: fmt.Sprintf("`%s`", coin.CoinData.CreatorAddr),
		},
		{
			Name:   "Dev Holding Amount",
			Value:  fmt.Sprintf("`%s`", convertDecimalToPercentage(coin.CoinData.DevHoldingAmount)),
			Inline: true,
		},
		{
			Name:   "Is Unique Coin",
			Value:  "not developed yet",
			Inline: true,
		},
		{
			Name:  "Socials",
			Value: socialsValue,
		},
		{
			Name:  "Links",
			Value: fmt.Sprintf("[SolScan](%s) | [PumpFun](%s)", "https://solscan.io/token/"+coin.CoinData.TokenAddr, "https://pump.fun/"+coin.CoinData.TokenAddr),
		},
	}
	return fields
}

func buildSocials(ipfsData *models.IPFS) string {
	if ipfsData == nil {
		return "N/A"
	}

	var socials []string
	if ipfsData.TelegramURL != "" {
		socials = append(socials, fmt.Sprintf("[Telegram](%s)", ipfsData.TelegramURL))
	}
	if ipfsData.TwitterURL != "" {
		socials = append(socials, fmt.Sprintf("[Twitter](%s)", ipfsData.TwitterURL))
	}
	if ipfsData.WebsiteURL != "" {
		socials = append(socials, fmt.Sprintf("[Website](%s)", ipfsData.WebsiteURL))
	}

	if len(socials) == 0 {
		return "N/A"
	}

	return strings.Join(socials, " | ")
}

func convertDecimalToPercentage(decimal float64) string {
	return fmt.Sprintf("%.2f%%", decimal*100)
}

func handleError(req *http.Response, body []byte) error {
	switch req.StatusCode {
	case http.StatusUnauthorized:
		logger.Log(logger.LevelError, "Webhook doesn't exist", logger.String("error", "Unauthorized"), logger.String("url", req.Request.URL.String()))
		return fmt.Errorf("unauthorized: %s", string(body))
	case http.StatusTooManyRequests:
		logger.Log(logger.LevelError, "Webhook rate limited", logger.String("error", "Rate limited"), logger.String("url", req.Request.URL.String()))
		return fmt.Errorf("rate limited: %s", string(body))
	default:
		logger.Log(logger.LevelError, fmt.Sprintf("Error sending webhook, code: %d", req.StatusCode), logger.String("error", string(body)), logger.String("url", req.Request.URL.String()))
		return fmt.Errorf("error sending webhook, code: %d, body: %s", req.StatusCode, string(body))
	}
}
