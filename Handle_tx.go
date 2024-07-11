package main

import (

)

//Declaing the coin struct 
type Coin struct {
}


var(
	//This is where we are gonna wanna set our enviroment variables
)

func handleTx(signature string) {
	// handle the tx
}

func getTxinfo(signature string) {
	// get tx info
}

func is_error(response string) {
	// check if there is an error
}

func get_mint_instruction(instruction string) {
	// get mint instruction
}

func is_create_tx(instruction string) {
	// check if it is a create tx
}

func parse_tx(tx string) {
	// parse the tx, we gonna wanna re write this function up as there is a comment that says clean it
}

func fetch_ipfs_and_update_coin_info(ipfs_url string) {
	// theres a note that this shouldnt be named like this but we can change it later maybe embed it into handle tx
}

