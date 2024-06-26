import websocket
import json
from handle_tx import handle_tx
### WHEN MOVING TO GO => HANLDE HTTP ERRORS GRACEFULLY + LOGGING 

# Log subscribe
"""
{
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
"""
# What will it manage: 
# 1) On open -> ws.send the log subscribe message
# 2) On message -> ws.recv the log message (this is where you'll 
#                   send the succesful log messages to the get tx function)
# 2.1) If the log message has a error -> skip it
# 3) on error -> handle the error gracefully
# 4) on close -> handle the close gracefully


#send request to get tx if no error -> check if there’s an instruction with 14 accounts -> decrypt instruction data’s first 8bytes  -> if matches create return webhook with name, symbol and ipfs url

from settings import NODE_URL

def on_message(ws, message):
    print(message)
    # handle tx if no error 
        # on a new thread
    handle_tx(message) # this will be a new thread

def on_error(ws, error):
    print(error)

def on_close(ws):
    print("### closed ###")

def on_open(ws):
    print("### open ###")
    # SEND MESSAGE HERE

if __name__ == "__main__":
    ws = websocket.WebSocketApp("", 
                                on_message = on_message,
                                on_error = on_error,
                                on_close = on_close,
                                on_open = on_open)
    ws.run_forever()