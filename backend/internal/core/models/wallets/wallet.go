package wallets

import "github.com/gagliardetto/solana-go"

type SolanaWallet struct {
	Id         string
	WalletName string
	PrivateKey solana.PrivateKey
	PublicKey  solana.PublicKey
}
