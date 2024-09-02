package transactions

import (
	"pump_fun/internal/models"

	"github.com/gagliardetto/solana-go"
)

func DecodeCreateTransaction(transaction *solana.Transaction) models.DecodedInstruction {
	i0 := transaction.Message.Instructions[3]
	decodedInstruction := DecodeInstruction(i0, transaction)
	decodedInstructionStruct := mapToStruct(decodedInstruction)
	return decodedInstructionStruct
}

func mapToStruct(decodedInstruction interface{}) models.DecodedInstruction {
	decodedInstructionMap, ok := decodedInstruction.(map[string]string)
	if !ok {
		panic("decodedInstruction is not of type map[string]string")
	}

	return models.DecodedInstruction{
		Name:     decodedInstructionMap["Name"],
		Symbol:   decodedInstructionMap["Symbol"],
		IPFS_URL: decodedInstructionMap["Uri"],
	}
}
