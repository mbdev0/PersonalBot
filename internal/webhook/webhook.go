package webhook

import (
	"pump_fun/models"
)

//We are gonna need to implement the same coin struct as in the get tx function, maybe find a package that allows models?

func FormatCoinInfo(coin models.Coin) string {
	//Format the coin info
	return ""
}

// copied params from the py implementation adjust them as needed
func SendTelegramMessage(botToken, chatID, message string) {
	//Send a message to the telegram bot
}

func SendDiscordMessage(webhookURL, message string) {
	//send a essage to the discord webhook
}

func PostRequest(url string, data interface{}) {
	// Send a post request to the url with the data
}
