package transactions

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"pump_fun/internal/constants"
	"pump_fun/internal/models"
	solanaclient "pump_fun/internal/solana-client"

	bin "github.com/gagliardetto/binary"

	"pump_fun/internal/logger"
)

func GetTransaction(signature string) (*solana.Transaction, error) {
	programID := solana.MustPublicKeyFromBase58(constants.ProgramID)
	solana.RegisterInstructionDecoder(programID, CustomInstructionDecoder)
	transaction, err := getTransaction(signature)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func getTransaction(signature string) (*solana.Transaction, error) {
	rpcClient := solanaclient.NewHttpClient()

	version := uint64(0)
	ctxTimeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	out, err := rpcClient.GetTransaction(
		ctx,
		solana.MustSignatureFromBase58(signature),
		&rpc.GetTransactionOpts{
			MaxSupportedTransactionVersion: &version,
			Encoding:                       solana.EncodingBase64,
		},
	)

	if err != nil {

		logger.Log(logger.LevelError, "Error getting transaction", logger.Error(err))
		return nil, err
	}

	transaction, err := solana.TransactionFromDecoder(bin.NewBinDecoder(out.Transaction.GetBinary()))
	if err != nil {
		logger.Log(logger.LevelError, "Error decoding transaction", logger.Error(err))
		return nil, err
	}
	return transaction, nil
}

func DecodeInstruction(compiledInstruction solana.CompiledInstruction, transaction *solana.Transaction) interface{} {
	progKey, err := transaction.ResolveProgramIDIndex(compiledInstruction.ProgramIDIndex)
	if err != nil {
		logger.Log(logger.LevelError, "Error decoding program ID", logger.Error(err))
	}

	accounts, err := compiledInstruction.ResolveInstructionAccounts(&transaction.Message)
	if err != nil {
		logger.Log(logger.LevelError, "Error decoding Instruction Accounts", logger.Error(err))
	}

	decodedInstruction, err := solana.DecodeInstruction(
		progKey,
		accounts,
		compiledInstruction.Data,
	)
	if err != nil {
		logger.Log(logger.LevelError, "Error decoding Instructions", logger.Error(err))
	}

	return decodedInstruction
}

// TODO: Remove DecodedInstruction and update this to be parse to MintData
func ParseTransaction(transaction *solana.Transaction) models.DecodedInstruction {
	compiledInstruction := transaction.Message.Instructions[3]
	decodedInstruction := DecodeInstruction(compiledInstruction, transaction)
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
		IPFS_URL: decodedInstructionMap["IPFS_URL"],
	}
}
