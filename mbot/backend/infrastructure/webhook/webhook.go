package webhook

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"encoding/json"
	"net/http"

	"personal_bot/backend/pkg/logger"
)

var (
	connectionTimeout = 30 * time.Second
)

func SendWebhook(discordWebhookURL string, message Webhook) error {
	return sendDiscordMessage(discordWebhookURL, message)
}

func sendDiscordMessage(webhookURL string, message Webhook) error {
	transport := &http.Transport{
		IdleConnTimeout: connectionTimeout,
	}
	client := &http.Client{Transport: transport}

	marshaledWebhook, err := json.Marshal(message)

	if err != nil {
		return err
	}

	req, err := client.Post(webhookURL, "application/json", bytes.NewBuffer(marshaledWebhook))

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(req.Body)

	if req.StatusCode != 204 {

		body, err := io.ReadAll(req.Body)

		if err != nil {
			return err
		}

		if err := handleError(req, body); err != nil {
			return err
		}
	}

	return nil
}

func handleError(req *http.Response, body []byte) error {
	switch req.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized: %s", string(body))
	case http.StatusTooManyRequests:
		return fmt.Errorf("rate limited: %s", string(body))
	default:
		return fmt.Errorf("error sending webhook, code: %d, body: %s", req.StatusCode, string(body))
	}
}
