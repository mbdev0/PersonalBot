package lookuptable

import (
	"context"
	"sync"

	"pump_fun/internal/constants"
	"pump_fun/internal/logger"
	rpcclient "pump_fun/internal/rpc_client"

	"github.com/gagliardetto/solana-go"
	lookup "github.com/gagliardetto/solana-go/programs/address-lookup-table"
)

var (
	addressLookupTable *lookup.AddressLookupTableState
	lock               = &sync.Mutex{}
)

func GetAddressLookupTable() *lookup.AddressLookupTableState {
	return getLookupTable()
}

func getLookupTable() *lookup.AddressLookupTableState {
	if addressLookupTable == nil {
		lock.Lock()
		defer lock.Unlock()
		if addressLookupTable == nil {
			addressLookupTable = loadLookupTable()
		}
	}

	return addressLookupTable
}

func loadLookupTable() *lookup.AddressLookupTableState {
	lookupTable, err := lookup.GetAddressLookupTable(context.Background(), rpcclient.GetClient(), solana.MustPublicKeyFromBase58(constants.AddressLookupTableAccount))

	if err != nil {
		logger.Error("Error trying to get the address lookup table: ", err)
	}

	return lookupTable
}
