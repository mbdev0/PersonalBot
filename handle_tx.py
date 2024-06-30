#Send a post request to node url with the following body:
"""
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "getTransaction",
  "params": [
    "[ENTER SIGNATURE HERE]",
    {
      "encoding": "json",
      "maxSupportedTransactionVersion": 0
    }
  ]
}
"""
# we need to check instruction data for create tx -> its in transaction: result -> transaction -> message -> instructions  
# this is where all the instructions of the tx lies
# find the instruction with 14 accounts OR if there is a better way e.g. loop through all instructions and check if there is a create within the 8 bytes
# of the decrypted instruction data
# if the first 8 bytes match a create tx 
"""
var (
	// Creates the global state.
	Instruction_Initialize = ag_binary.TypeID([8]byte{175, 175, 109, 31, 13, 152, 155, 237})

	// Sets the global state parameters.
	Instruction_SetParams = ag_binary.TypeID([8]byte{27, 234, 178, 52, 147, 2, 187, 141})

	// Creates a new coin and bonding curve.
	Instruction_Create = ag_binary.TypeID([8]byte{24, 30, 200, 40, 5, 28, 7, 119})

	// Buys tokens from a bonding curve.
	Instruction_Buy = ag_binary.TypeID([8]byte{102, 6, 61, 18, 1, 218, 235, 234})

	// Sells tokens into a bonding curve.
	Instruction_Sell = ag_binary.TypeID([8]byte{51, 230, 133, 164, 1, 127, 131, 173})

	// Allows the admin to withdraw liquidity for a migration once the bonding curve completes
	Instruction_Withdraw = ag_binary.TypeID([8]byte{183, 18, 70, 156, 148, 109, 161, 34})
)
"""

# testing for parsing the instruction data + TX's when doing GO

from Models.Coin import Coin 
from settings import HTTP_NODE_URL
from webhook import send_telegram_webhook
import httpx
import base58

CREATE_TX_BYTES = [24, 30, 200, 40, 5, 28, 7, 119]

def handle_tx(signature: str):
	transaction = get_tx_info(signature)
	instruction_data = get_mint_instruction_data(transaction)
	coin = parse_tx(instruction_data)
	fetch_ifps_links_and_update_coin(coin)
	print(coin.dict())
	send_telegram_webhook(coin)
	
def get_tx_info(signature) -> dict:
	getTransaction_body = {
		"jsonrpc": "2.0",
		"id": 1,
		"method": "getTransaction",
		"params": [
			signature,
			{
				"encoding": "json",
				"maxSupportedTransactionVersion": 0
			}
		]
	}

	response = httpx.post(HTTP_NODE_URL, json=getTransaction_body)
	isErrorPresent = is_error_present(response.json())

	if response.status_code != 200 or isErrorPresent:
		raise Exception("Error in getting tx info")
	return response.json()

def is_error_present(response: dict) -> bool:
    if response.get("error") is not None:
        return True
    result = response.get("result")
    if result and result.get("meta") and result.get("meta").get("err") is not None:
        return True
    return False

def get_mint_instruction_data(instruction_data: dict) -> dict:
	program_instructions = instruction_data.get("result").get("transaction").get("message").get("instructions")
	for instruction in program_instructions:
		if is_create_tx(instruction):
			return instruction
	raise Exception("No mint instruction found")

def is_create_tx(instruction: dict) -> bool:
	if instruction.get("data") and instruction.get("accounts") and len(instruction.get("accounts")) == 14:
		data = instruction.get("data")
		decoded_data = base58.b58decode(data)
		byte_data = list(decoded_data)
		if byte_data[:8] == CREATE_TX_BYTES:
			return True
	return False

def parse_tx(instruction: dict) -> Coin:
	#TODO: clean this function up if possible

	data = instruction.get("data")
	decoded_data = base58.b58decode(data)
	byte_data = list(decoded_data)
	byte_data = byte_data[8:]
	
	byte_data = [byte for byte in byte_data if byte != 0]
	
	name_length = byte_data[0]
	name = byte_data[1:1+name_length]
	
	byte_data = byte_data[1+name_length:]
	
	symbol_length = byte_data[0]
	symbol = byte_data[1:1+symbol_length]
	
	byte_data = byte_data[1+symbol_length:]
	
	uri_length = byte_data[0]
	uri = byte_data[1:1+uri_length]

	mint_name = bytes(name) 
	mint_symbol = bytes(symbol)
	mint_uri = bytes(uri)

	coin = Coin(name=mint_name.decode('utf-8'), symbol=mint_symbol.decode('utf-8'), ipfs_url=mint_uri.decode('utf-8'))
	return coin

def fetch_ifps_links_and_update_coin(coin):
	with httpx.Client(follow_redirects=True) as client:
			response = client.get(coin.ipfs_url)
			data = response.json()

			# Update the coin object with the fetched data
			# TODO : Clean this up in go, shouldn't really be having side effects happen in the same function
			coin.telegram_url=data.get("telegram")
			coin.twitter_url=data.get("twitter")
			coin.website_url=data.get("website")
			coin.image_url=data.get("image")

# if __name__ == "__main__":
# 	handle_tx("5vQ5yPrGjE6LX5ZoLPCXK6YQtKymXLpeyVqBM6g2NUU86pvSSRVsTMG1FTyeTLPWErqSV8KAT2gD8bmEK6fzFVQg")