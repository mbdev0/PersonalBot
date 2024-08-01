package main

//We are gonna need to implement the same coin struct as in the get tx function, maybe find a package that allows models?

func formatCoinInfo(coin Coin) string {
	//Format the coin info
	return ""
}

// copied params from the py implementation adjust them as needed
func sendTelegramMessage(botToken, chatID, message string) {
	//Send a message to the telegram bot
}

func sendDiscordMessage(webhookURL, message string) {
	//send a essage to the discord webhook
}

func postRequest(url string, data interface{}) {
	// Send a post request to the url with the data
}
