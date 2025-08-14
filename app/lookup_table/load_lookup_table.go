package lookuptable

import (
	"context"
	"sync"

	"pump_fun/internal/core/constants"
	"pump_fun/internal/solana/client"
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
	lookupTable, err := lookup.GetAddressLookupTable(context.Background(), client.GetClient(), solana.MustPublicKeyFromBase58(constants.AddressLookupTableAccount))

	if err != nil {
		logger.Error("Error trying to get the address lookup table: ", err)
		return nil
	}

	accountLookupMap := make(map[solana.PublicKey]solana.PublicKeySlice)
	accountLookupTableAccount := solana.MustPublicKeyFromBase58(constants.AddressLookupTableAccount)
	accountLookupMap[accountLookupTableAccount] = lookupTable.Addresses

	return &accountLookupMap
}
