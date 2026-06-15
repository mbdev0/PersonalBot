package instructions

import (
	"context"
	"fmt"
	"personal_bot/backend/internal/core/constants"
	"personal_bot/backend/internal/solana/client"
	"personal_bot/backend/internal/solana/pda"
	"personal_bot/backend/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetIdempotentInstruction(ctx context.Context, wallet solana.PublicKey, mintAddress solana.PublicKey, httpNode *rpc.Client) (*solana.GenericInstruction, error) {
	idEmponentInstruction, err := getIdempotentInstructionIfExists(ctx, wallet, mintAddress, httpNode)

	if err != nil {
		logger.Error("Error getting IdempotentInstruction: ", err)
		return nil, err
	}

	return idEmponentInstruction, nil
}

func getIdempotentInstructionIfExists(ctx context.Context, wallet solana.PublicKey, mintAddress solana.PublicKey, httpNode *rpc.Client) (*solana.GenericInstruction, error) {
	var associatedTokenAddressPubkey solana.PublicKey
	var err error

	isNewTokenProgram, err := IsTokenAccountNew(ctx, mintAddress, httpNode)
	if err != nil {
		return nil, err
	}

	if isNewTokenProgram {
		associatedTokenAddressPubkey, _, err = pda.FindToken2022AssociatedTokenAddress(wallet, mintAddress)
		if err != nil {
			return nil, err
		}
	} else {
		associatedTokenAddressPubkey, _, err = pda.FindTokenAssociatedTokenAddress(wallet, mintAddress)
		if err != nil {
			return nil, err
		}
	}

	_, err = client.GetAccountInfo(ctx, associatedTokenAddressPubkey.String(), httpNode)

	if err != nil {
		if err.Error() == "not found" {
			ataInstruction, err := makeAssociatedTokenAccountInstruction(wallet, wallet, mintAddress, associatedTokenAddressPubkey, isNewTokenProgram)
			if err != nil {
				return nil, err
			}

			return ataInstruction, nil
		}
		return nil, err
	}

	return nil, nil
}

func makeAssociatedTokenAccountInstruction(payer solana.PublicKey, walletAddress solana.PublicKey, splTokenAddress solana.PublicKey, associatedTokenAddressPubkey solana.PublicKey, isNewTokenProgram bool) (*solana.GenericInstruction, error) {
	associatedTokenAccountProgram := solana.MustPublicKeyFromBase58(constants.AssociatedTokenProgram)

	var tokenAccount solana.PublicKey
	if isNewTokenProgram {
		tokenAccount = solana.Token2022ProgramID
	} else {
		tokenAccount = solana.TokenProgramID
	}

	// get accounts
	accounts := []*solana.AccountMeta{
		solana.NewAccountMeta(payer, true, true),
		solana.NewAccountMeta(associatedTokenAddressPubkey, true, false),
		solana.NewAccountMeta(walletAddress, true, true),
		solana.NewAccountMeta(splTokenAddress, false, false),
		solana.NewAccountMeta(solana.SystemProgramID, false, false),
		solana.NewAccountMeta(tokenAccount, false, false),
	}

	// then make instruction
	inst := solana.NewInstruction(
		associatedTokenAccountProgram,
		accounts,
		[]byte{0},
	)

	return inst, nil
}

func IsTokenAccountNew(ctx context.Context, mintAddress solana.PublicKey, httpNode *rpc.Client) (bool, error) {
	accountInfo, err := client.GetAccountInfo(ctx, mintAddress.String(), httpNode)
	if err != nil {
		return false, err
	}

	if accountInfo.Value == nil {
		return false, fmt.Errorf("error whilst get account info on mint address")
	}

	ownerAccount := accountInfo.Value.Owner

	switch ownerAccount {
	case solana.TokenProgramID:
		return false, nil
	case solana.Token2022ProgramID:
		return true, nil
	}

	return false, fmt.Errorf("error whilst finding if token account is new")

}
