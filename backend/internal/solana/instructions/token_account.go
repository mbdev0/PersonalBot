package instructions

import (
	"context"
	"fmt"
	"pump_fun/internal/core/constants"
	"pump_fun/internal/solana/client"
	"pump_fun/internal/solana/programs/pumpfun/pda"
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go"
)

func GetIdempotentInstruction(wallet solana.PublicKey, mintAddress solana.PublicKey, ctx context.Context) (*solana.GenericInstruction, error) {
	idEmponentInstruction, err := getIdempotentInstructionIfExists(wallet, mintAddress, ctx)

	if err != nil {
		logger.Error("Error getting IdempotentInstruction: ", err)
		return nil, err
	}

	return idEmponentInstruction, nil
}

func getIdempotentInstructionIfExists(wallet solana.PublicKey, mintAddress solana.PublicKey, ctx context.Context) (*solana.GenericInstruction, error) {
	var associatedTokenAddressPubkey solana.PublicKey
	var err error

	isNewTokenProgram, err := IsTokenAccountNew(mintAddress, ctx)
	if err != nil {
		return nil, err
	}

	if isNewTokenProgram {
		associatedTokenAddressPubkey, _, err = pda.FindToken2022AssociatedTokenAddress(wallet, mintAddress)
	} else {
		associatedTokenAddressPubkey, _, err = pda.FindTokenAssociatedTokenAddress(wallet, mintAddress)
	}

	_, err = client.GetAccountInfo(associatedTokenAddressPubkey.String(), ctx)

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

func IsTokenAccountNew(mintAddress solana.PublicKey, ctx context.Context) (bool, error) {
	accountInfo, err := client.GetAccountInfo(mintAddress.String(), ctx)
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
