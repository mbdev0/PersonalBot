package client

import (
	"context"
	"fmt"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/solana/pda"
	"strconv"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

var (
	TokenProgramV1 = solana.MustPublicKeyFromBase58(constants.TokenProgram)
	TokenProgramV2 = solana.MustPublicKeyFromBase58(constants.Token2022Program)
)

func GetTokenAccountBalance(ctx context.Context, associatedTokenAddress solana.PublicKey, rpcClient *rpc.Client) (tokenAmount *uint64, err error) {

	result, err := rpcClient.GetTokenAccountBalance(ctx, associatedTokenAddress, rpc.CommitmentConfirmed)
	if err != nil {
		return nil, err
	}

	amount, err := strconv.ParseInt(result.Value.Amount, 10, 64)
	if err != nil {
		return nil, err
	}

	uintAmount := uint64(amount)

	return &uintAmount, nil
}

func GetTokenProgramForMint(ctx context.Context, mint solana.PublicKey, rpcClient *rpc.Client) (solana.PublicKey, error) {

	accountInfo, err := rpcClient.GetAccountInfo(ctx, mint)
	if err != nil {
		return solana.PublicKey{}, err
	}
	if accountInfo == nil || accountInfo.Value == nil {
		return solana.PublicKey{}, fmt.Errorf("mint account not found: %s", mint)
	}

	owner := accountInfo.Value.Owner
	switch owner {
	case TokenProgramV1, TokenProgramV2:
		return owner, nil
	default:
		return solana.PublicKey{}, fmt.Errorf("unknown token program owner: %s", owner)
	}
}

func GetATA(ctx context.Context, wallet solana.PublicKey, mint solana.PublicKey, rpcClient *rpc.Client) (solana.PublicKey, error) {
	tokenProgram, err := GetTokenProgramForMint(ctx, mint, rpcClient)
	if err != nil {
		return solana.PublicKey{}, err
	}

	var ata solana.PublicKey
	switch tokenProgram {
	case solana.Token2022ProgramID:
		ata, _, err = pda.FindToken2022AssociatedTokenAddress(wallet, mint)
	default:
		ata, _, err = pda.FindTokenAssociatedTokenAddress(wallet, mint)
	}

	return ata, err
}
