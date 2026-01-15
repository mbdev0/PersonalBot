package mapper

import (
	"fmt"
	"personal_bot/infrastructure/persistence/models"
	"personal_bot/internal/core/models/wallets"

	"github.com/gagliardetto/solana-go"
)

func WalletRepoToWallet(src models.WalletRepository) (wallets.SolanaWallet, error) {

	privateKey, err := solana.PrivateKeyFromBase58(src.PrivateKey)
	if err != nil {
		return wallets.SolanaWallet{}, fmt.Errorf("error whilst mapping: %v", err)
	}

	publicKey := privateKey.PublicKey()

	return wallets.SolanaWallet{
		Id:         src.Id,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		WalletName: src.WalletName,
	}, nil
}

func WalletToWalletRepo(src wallets.SolanaWallet) models.WalletRepository {
	return models.WalletRepository{
		Id:         src.Id,
		WalletName: src.WalletName,
		PrivateKey: src.PrivateKey.String(),
		Chain:      "Solana",
	}
}
