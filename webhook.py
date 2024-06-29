
from Models.Coin import Coin
from settings import BOT_TOKEN, CHAT_ID
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

# test_coin = {'name': '1st Blink Sydney Sweeney', 'symbol': 'blinkboobs', 'ipfs_url': 'https://cf-ipfs.com/ipfs/Qme9GvUAU53DZinju8areGH9cA9KsA1a6YRTKCpbeMsUkt', 'telegram_url': None, 'twitter_url': None, 'website_url': None}
# send_telegram_webhook(test_coin)