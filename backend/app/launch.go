package app

import (
	lookuptable "personal_bot/app/lookup_table"
	"personal_bot/app/pumpfun_idl"
	"personal_bot/app/validator"
	"personal_bot/infrastructure/solana_price"
	"personal_bot/pkg/logger"

	"github.com/gagliardetto/solana-go/rpc"
)

func Launch(generalNode *rpc.Client) {
	pumpfun_idl.GetIdlMap()
	validator.GetValidator()
	if _, err := lookuptable.GetAddressLookupTable(generalNode); err != nil {
		logger.Error(err)
	}

	_, err := solana_price.GetSolPrice()
	if err != nil {
		return
	}
}
