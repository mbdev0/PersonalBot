package app

import (
	lookuptable "personal_bot/app/lookup_table"
	"personal_bot/app/pumpfun_idl"
	"personal_bot/app/validator"
	"personal_bot/infrastructure/solana_price"

	"github.com/gagliardetto/solana-go/rpc"
)

func Launch(generalNode *rpc.Client) {
	pumpfun_idl.GetIdlMap()
	_, err := solana_price.GetSolPrice()
	if err != nil {
		return
	}
	lookuptable.GetAddressLookupTable(generalNode)
	validator.GetValidator()
}
