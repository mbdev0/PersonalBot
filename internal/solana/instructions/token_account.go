package instructions

import (
	"pump_fun/internal/core/models"
	"pump_fun/internal/solana/client"
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
)

func GetIdempotentInstruction(wallet solana.PublicKey, mintAddress solana.PublicKey, cancellationToken models.CancelToken) *associatedtokenaccount.Instruction {
	idEmponentInstruction, err := getIdempotentInstructionIfExists(wallet, mintAddress, cancellationToken)

	if err != nil {
		logger.Error("Error getting IdempotentInstruction: ", err)
		return nil
	}

	return idEmponentInstruction
}

func getIdempotentInstructionIfExists(wallet solana.PublicKey, mintAddress solana.PublicKey, cancellationToken models.CancelToken) (*associatedtokenaccount.Instruction, error) {
	associatedTokenAddressPubkey, _, err := solana.FindAssociatedTokenAddress(wallet, mintAddress)
	if err != nil {
		return nil, err
	}

	_, err = client.GetAccountInfo(associatedTokenAddressPubkey.String(), cancellationToken)

	if err != nil {
		if err.Error() == "not found" {
			ataInstruction := associatedtokenaccount.NewCreateInstruction(wallet, wallet, mintAddress).Build()
			return ataInstruction, nil
		}
		return nil, err
	}

	return nil, nil
}
