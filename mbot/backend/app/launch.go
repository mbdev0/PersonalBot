package app

import (
	lookuptable "personal_bot/backend/app/lookup_table"
	"personal_bot/backend/app/pumpfun_idl"
	"personal_bot/backend/app/validator"
	"personal_bot/backend/infrastructure/solana_price"
	"personal_bot/backend/pkg/logger"

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
