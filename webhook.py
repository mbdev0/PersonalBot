
from Models.Coin import Coin
from settings import BOT_TOKEN, CHAT_ID, DISCORD_WEBHOOK
from discord_webhook import DiscordWebhook, DiscordEmbed
import threading
import requests

def send_telegram_webhook(coin: Coin):
    coin_info = format_coint_info(coin.dict())
    send_telegram_message(BOT_TOKEN, CHAT_ID, coin_info)

def format_coint_info(coin: dict) -> str:
    # Need to be updated to a better format
    return "\n".join([f"{key}: {value if value is not None else 'N/A'}" for key, value in coin.items()])

def send_telegram_message(bot_token: str, chat_id: str, message: str):
    url = f"https://api.telegram.org/bot{bot_token}/sendMessage"
    data = {
        "chat_id": chat_id,
        "text": message
    }
    response = requests.post(url, data=data)
    print("Sent the following to telegram successfully:\n"+message)

def send_discord_webhook(coin: Coin):
    coin_info = format_coint_info(coin.dict())
    send_discord_message(DISCORD_WEBHOOK, coin_info)

def send_discord_message(webhook_url: str, message: str):
    webhook = DiscordWebhook(url=webhook_url)
    embed = DiscordEmbed(title='New Coin Alert!', description=message, color=242424)
    webhook.add_embed(embed)
    webhook.execute()
    print("Sent the following to discord successfully:\n"+message)



# if __name__ == "__main__":
    
#     test_coin = {'name': '1st Blink Sydney Sweeney', 'symbol': 'blinkboobs', 'ipfs_url': 'https://cf-ipfs.com/ipfs/Qme9GvUAU53DZinju8areGH9cA9KsA1a6YRTKCpbeMsUkt', 'telegram_url': None, 'twitter_url': None, 'website_url': None}
#     test_coin = Coin(**test_coin)

#     # send both webhook at the same time using threading
#     threading.Thread(target=send_discord_webhook, args=(test_coin,)).start()
#     threading.Thread(target=send_telegram_webhook, args=(test_coin,)).start()