package notifier

import (
	"fmt"
	"personal_bot/infrastructure/webhook"
	"personal_bot/internal/core/notifier"
	"personal_bot/internal/core/settings"
)

type DiscordNotifier struct {
	DiscordWebhook string
	SendOnFail     bool
	SendOnSuccess  bool
}

func NewDiscordNotifier(settings settings.Settings) *DiscordNotifier {
	return &DiscordNotifier{
		DiscordWebhook: settings.DiscordWebhook,
		SendOnFail:     settings.SendOnFail,
		SendOnSuccess:  settings.SendOnSuccess,
	}
}

func (dn *DiscordNotifier) Update(settings settings.Settings) {
	dn.DiscordWebhook = settings.DiscordWebhook
	dn.SendOnFail = settings.SendOnFail
	dn.SendOnSuccess = settings.SendOnSuccess
}

// simple test webhook send
func (dn *DiscordNotifier) TestDiscordEndpoint(discordWebhookUrl string) error {
	testWebhook := dn.setupTestMessage()
	err := webhook.SendWebhook(discordWebhookUrl, testWebhook)
	if err != nil {
		return err
	}
	return nil
}

func (dn *DiscordNotifier) SendFailure(payload notifier.ErrorNotifierPayload) error {
	if dn.DiscordWebhook == "" || !dn.SendOnFail {
		return nil
	}
	return webhook.SendWebhook(dn.DiscordWebhook, dn.setupFailedMessage(payload))
}

func (dn *DiscordNotifier) SendSuccessBuy(payload notifier.BuyNotifierPayload) error {
	if dn.DiscordWebhook == "" || !dn.SendOnSuccess {
		return nil
	}
	return webhook.SendWebhook(dn.DiscordWebhook, dn.setupSuccessBuyMessage(payload))

}

func (dn *DiscordNotifier) SendSuccessSell(payload notifier.SellNotifierPayload) error {
	if dn.DiscordWebhook == "" || !dn.SendOnSuccess {
		return nil
	}
	return webhook.SendWebhook(dn.DiscordWebhook, dn.setupSuccessSellMessage(payload))

}

func (dn *DiscordNotifier) setupTestMessage() webhook.Webhook {
	embed := webhook.Embeds{
		Title: "SUCCESS - Testing Webhook",
		Color: 5814783,
		Author: webhook.Author{
			Name: "AIO BOT - Test Webhook",
		},
		Footer: webhook.Footer{
			Text: "AIO Bot - Mustafa",
		},
	}

	testWebhook := webhook.Webhook{
		Embeds: []webhook.Embeds{embed},
	}
	return testWebhook
}

func (dn *DiscordNotifier) setupSuccessBuyMessage(payload notifier.BuyNotifierPayload) webhook.Webhook {
	var title string
	switch payload.TaskType {
	case "BUY":
		title = fmt.Sprintf("[QB] - %s", payload.TokenAddress)
	default:
		title = fmt.Sprintf("%s %d - Task %s [BUY]", payload.TaskType, *payload.StrategyId, shortAddress(payload.TokenAddress))
	}

	embed := webhook.Embeds{
		Title: title,
		Color: 579080,
		Author: webhook.Author{
			Name: "AIO BOT - Success",
		},
		Fields: []webhook.Fields{
			{
				Name:   "Amount Bought (SOL)",
				Value:  fmt.Sprintf("`%f`", payload.AmountPaid),
				Inline: true,
			},
			{
				Name:   "Tokens",
				Value:  fmt.Sprintf("`%f`", payload.TokensBought),
				Inline: true,
			},
			{
				Name:  "Wallet Address",
				Value: fmt.Sprintf("`%s`", payload.WalletAddress),
			},
			{
				Name:  "Token Address",
				Value: fmt.Sprintf("`%s`", payload.TokenAddress),
			},
			{
				Name: "Links",
				//Solscan, Axiom, Pumpfun
				Value: fmt.Sprintf("[SolScan](%s) | [Axiom](%s) | [PumpFun](%s)", "https://solscan.io/tx/"+payload.TxHash, "https://axiom.trade/meme/"+payload.BondingCurve+"?chain=sol", "https://pump.fun/"+payload.TokenAddress),
			},
		},
		Footer: webhook.Footer{
			Text: "AIO Bot - Mustafa",
		}}

	return webhook.Webhook{Embeds: []webhook.Embeds{embed}}
}

func (dn *DiscordNotifier) setupSuccessSellMessage(payload notifier.SellNotifierPayload) webhook.Webhook {
	var title string
	switch payload.TaskType {
	case "SELL":
		title = fmt.Sprintf("[QS] - %s", payload.TokenAddress)
	case "SellTask":
		title = fmt.Sprintf("[SELL] - %s", payload.TokenAddress)
	default:
		title = fmt.Sprintf("%s %d - Task %s [SELL]", payload.TaskType, *payload.StrategyId, shortAddress(payload.TokenAddress))
	}

	embed := webhook.Embeds{
		Title: title,
		Color: 579080,
		Author: webhook.Author{
			Name: "AIO BOT - Success",
		},
		Fields: []webhook.Fields{
			{
				Name:   "Recieved from Sale (SOL)",
				Value:  fmt.Sprintf("`%f`", payload.AmountSold),
				Inline: true,
			},
			{
				Name:   "Tokens Sold",
				Value:  fmt.Sprintf("`%f`", payload.TokensSold),
				Inline: true,
			},
			{Name: "\u200b", Value: "\u200b", Inline: true},
			{
				Name:   "Tokens Remaining",
				Value:  fmt.Sprintf("`%f`", payload.TokensRemaining),
				Inline: true,
			},
			{
				Name:   "Current Profit",
				Value:  fmt.Sprintf("`%f`", payload.CurrentProfit),
				Inline: true,
			},
			{Name: "\u200b", Value: "\u200b", Inline: true},
			{
				Name:  "Wallet Address",
				Value: fmt.Sprintf("`%s`", payload.WalletAddress),
			},
			{
				Name: "Links",
				//Solscan, Axiom, Pumpfun
				Value: fmt.Sprintf("[SolScan](%s) | [Axiom](%s) | [PumpFun](%s)", "https://solscan.io/tx/"+payload.TxHash, "https://axiom.trade/meme/"+payload.BondingCurve+"?chain=sol", "https://pump.fun/"+payload.TokenAddress),
			},
		},
		Footer: webhook.Footer{
			Text: "AIO Bot - Mustafa",
		}}

	return webhook.Webhook{Embeds: []webhook.Embeds{embed}}
}

func (dn *DiscordNotifier) setupFailedMessage(payload notifier.ErrorNotifierPayload) webhook.Webhook {
	var title string
	switch payload.TaskType {
	case "BUY":
		title = fmt.Sprintf("[QB] - %s", payload.TokenAddress)
	case "SELL":
		title = fmt.Sprintf("[QS] - %s", payload.TokenAddress)
	default:
		if payload.StrategyId == nil {
			title = fmt.Sprintf("%s - Task %s", payload.TaskType, shortAddress(payload.TokenAddress))
		} else {
			title = fmt.Sprintf("%s %d - Task %s", payload.TaskType, *payload.StrategyId, shortAddress(payload.TokenAddress))
		}
	}

	embed := webhook.Embeds{
		Title: title,
		Color: 15869735,
		Author: webhook.Author{
			Name: "AIO BOT - Failure",
		},
		Fields: []webhook.Fields{
			{
				Name:   "Failure Reason",
				Value:  fmt.Sprintf("`%s`", payload.Error),
				Inline: true,
			},
			{
				Name:  "Wallet Address",
				Value: fmt.Sprintf("`%s`", payload.WalletAddress),
			},
			{
				Name: "Links",
				//Solscan, Axiom, Pumpfun
				Value: fmt.Sprintf("[SolScan](%s) | [Axiom](%s) | [PumpFun](%s)", "https://solscan.io/tx/"+payload.TxHash, "https://axiom.trade/meme/"+payload.BondingCurve+"?chain=sol", "https://pump.fun/"+payload.TokenAddress),
			},
		},
		Footer: webhook.Footer{
			Text: "AIO Bot - Mustafa",
		}}

	return webhook.Webhook{Embeds: []webhook.Embeds{embed}}
}

func shortAddress(s string) string {
	if len(s) <= 8 {
		return s
	}
	return s[:4] + "..." + s[len(s)-4:]
}
