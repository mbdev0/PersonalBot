package lookuptable

import (
	"context"
	"sync"

	"personal_bot/internal/core/constants"
	"personal_bot/pkg/logger"

	"github.com/gagliardetto/solana-go"
	lookup "github.com/gagliardetto/solana-go/programs/address-lookup-table"
	"github.com/gagliardetto/solana-go/rpc"
)

var (
	addressLookupTable map[solana.PublicKey]solana.PublicKeySlice
	mu                 sync.Mutex
)

func GetAddressLookupTable(node *rpc.Client) (map[solana.PublicKey]solana.PublicKeySlice, error) {
	mu.Lock()
	defer mu.Unlock()

	if addressLookupTable != nil {
		return addressLookupTable, nil
	}

	table, err := loadLookupTable(node)
	if err != nil {
		return nil, err
	}

	addressLookupTable = table
	return addressLookupTable, nil
}

func loadLookupTable(node *rpc.Client) (map[solana.PublicKey]solana.PublicKeySlice, error) {
	lookupTable, err := lookup.GetAddressLookupTable(
		context.Background(),
		node,
		solana.MustPublicKeyFromBase58(constants.AddressLookupTableAccount),
	)
	if err != nil {
		logger.Error("Error trying to get the address lookup table: ", err)
		return nil, err
	}

	accountLookupMap := make(map[solana.PublicKey]solana.PublicKeySlice)
	accountLookupMap[solana.MustPublicKeyFromBase58(constants.AddressLookupTableAccount)] = lookupTable.Addresses

	return accountLookupMap, nil
}
