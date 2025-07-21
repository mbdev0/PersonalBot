package lookuptable

import (
	"context"
	"sync"

	"pump_fun/internal/constants"
	rpcclient "pump_fun/internal/rpc_client"
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go"
	lookup "github.com/gagliardetto/solana-go/programs/address-lookup-table"
)

var (
	addressLookupTable *map[solana.PublicKey]solana.PublicKeySlice
	lock               = &sync.Mutex{}
)

func GetAddressLookupTable() map[solana.PublicKey]solana.PublicKeySlice {
	return getLookupTable()
}

func getLookupTable() map[solana.PublicKey]solana.PublicKeySlice {
	if addressLookupTable == nil {
		lock.Lock()
		defer lock.Unlock()
		if addressLookupTable == nil {
			addressLookupTable = loadLookupTable()
		}
	}

	return *addressLookupTable
}

func loadLookupTable() *map[solana.PublicKey]solana.PublicKeySlice {
	lookupTable, err := lookup.GetAddressLookupTable(context.Background(), rpcclient.GetClient(), solana.MustPublicKeyFromBase58(constants.AddressLookupTableAccount))

	if err != nil {
		logger.Error("Error trying to get the address lookup table: ", err)
	}

	accountLookupMap := make(map[solana.PublicKey]solana.PublicKeySlice)
	accountLookupTableAccount := solana.MustPublicKeyFromBase58(constants.AddressLookupTableAccount)
	accountLookupMap[accountLookupTableAccount] = lookupTable.Addresses

	return &accountLookupMap
}
