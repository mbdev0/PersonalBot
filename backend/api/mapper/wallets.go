package mapper

import (
	"personal_bot/api/dto"
	"personal_bot/internal/core/wallets"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
)

func MapWalletDtoToWallet(src dto.RequestWalletDto) (wallets.SolanaWallet, error) {
	id := uuid.NewString()
	privateKey, err := solana.PrivateKeyFromBase58(src.Private_key)
	if err != nil {
		return wallets.SolanaWallet{}, err
	}
	publicKey := privateKey.PublicKey()

	return wallets.SolanaWallet{
		Id:         id,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		WalletName: src.Wallet_name,
	}, nil
}

func MapWalletToDto(src wallets.SolanaWallet) dto.ResponseWalletDto {
	return dto.ResponseWalletDto{
		Id:         src.Id,
		WalletName: src.WalletName,
		PublicKey:  src.PrivateKey.PublicKey().String(),
		Chain:      "Solana",
	}
}
