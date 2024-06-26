import websocket 
from settings import WEBSOCKET_NODE_URL
import json
from handle_tx import handle_tx
### WHEN MOVING TO GO => HANLDE HTTP ERRORS GRACEFULLY + LOGGING 

# Log subscribe

# What will it manage: 
# 1) On open -> ws.send the log subscribe message
# 2) On message -> ws.recv the log message (this is where you'll 
#                   send the succesful log messages to the get tx function)
# 2.1) If the log message has a error -> skip it
# 3) on error -> handle the error gracefully
# 4) on close -> handle the close gracefully


#send request to get tx if no error -> check if there’s an instruction with 14 accounts -> decrypt instruction data’s first 8bytes  -> if matches create return webhook with name, symbol and ipfs url

from settings import WEBSOCKET_NODE_URL


def on_message(ws, message):
    try:
        data = json.loads(message)  # Try and parse the message as JSON
        if data["method"] == "logsNotification":
            result = data["params"]["result"]
            #print("Result:", result)  # debugging 
            if result["value"]["err"] is None:  # Little error check
                signature = result["value"]["signature"]
                try:
                    handle_tx(signature)  # Musty said only the signature is needed
                except Exception as e:
                    print("Error in handle_tx:", e)
            else:
                print("Transaction has an error, skipping...")
    except Exception as e:
        print("Response isn't in the expected format:", e)

def on_error(ws, error):
    print(error)

def on_close(ws):
    print("### closed ###")

def on_open(ws):
    print("WebSocket opened")
    subscribe_message = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "logsSubscribe",
        "params": [
            {
                "mentions": ["6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P"]
            },
            {
                "commitment": "finalized"
            }
        ]
    }
    ws.send(json.dumps(subscribe_message))

if __name__ == "__main__":
    ws = websocket.WebSocketApp(WEBSOCKET_NODE_URL, 
                                on_message = on_message,
                                on_error = on_error,
                                on_close = on_close,
                                on_open = on_open)
    ws.run_forever()