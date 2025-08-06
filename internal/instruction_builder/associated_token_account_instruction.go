package instructionbuilder

import (
	"pump_fun/internal/models"
	rpcclient "pump_fun/internal/rpc_client"
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
)

func GetIdEmponentInstruction(wallet solana.PublicKey, mintAddress solana.PublicKey, cancellationToken models.CancelToken) *associatedtokenaccount.Instruction {
	idEmponentInstruction, err := getIdEmponentInstructionIfExists(wallet, mintAddress, cancellationToken)

	if err != nil {
		logger.Error("Error getting IdEmponentInstruction: ", err)
		return nil
	}

	return idEmponentInstruction
}

func getIdEmponentInstructionIfExists(wallet solana.PublicKey, mintAddress solana.PublicKey, cancellationToken models.CancelToken) (*associatedtokenaccount.Instruction, error) {
	associatedTokenAddressPubkey, _, err := solana.FindAssociatedTokenAddress(wallet, mintAddress)
	if err != nil {
		return nil, err
	}

	_, err = rpcclient.GetAccountInfo(associatedTokenAddressPubkey.String(), cancellationToken)

	if err != nil {
		if err.Error() == "not found" {
			ataInstruction := associatedtokenaccount.NewCreateInstruction(wallet, wallet, mintAddress).Build()
			return ataInstruction, nil
		}
		return nil, err
	}

	return nil, nil
}
