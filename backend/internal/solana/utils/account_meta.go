package utils

import "github.com/gagliardetto/solana-go"

func GetAccountMeta(pubkey string, isMut bool, isSign bool) *solana.AccountMeta {
	return solana.NewAccountMeta(solana.MustPublicKeyFromBase58(pubkey), isMut, isSign)
}
