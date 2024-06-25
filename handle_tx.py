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

def handle_tx(signature: str):
	pass
	transaction = get_tx_info(signature)
	# check if there is a create tx
	if is_create_tx(transaction):
		coin = parse_tx(transaction)
		# send webhook with coin info
		# send_webhook(coin)
		pass
	
def get_tx_info(signature) -> str:
	# send a post req to node with the body template above
	pass

def is_create_tx(instruction_data: str) -> bool:
	# CHECK IF 14 ACCOUNTS or if better way try it
	pass

def parse_tx(transaction: str) -> Coin:
	pass # throws 
