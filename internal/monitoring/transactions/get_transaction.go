package transactions

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	solanaclient "pump_fun/internal/solana-client"

	bin "github.com/gagliardetto/binary"

	"pump_fun/internal/logger"
	"pump_fun/internal/models"
)

func GetTransaction(signature string) (models.DecodedInstruction, error) {
	programID := solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")
	solana.RegisterInstructionDecoder(programID, CustomInstructionDecoder)
	transaction := getTransaction(signature)
	return decodeCreateTransaction(transaction)
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

func decodeCreateTransaction(transaction *solana.Transaction) models.DecodedInstruction {
	i0 := transaction.Message.Instructions[3]
	progKey, err := transaction.ResolveProgramIDIndex(i0.ProgramIDIndex)
	if err != nil {
		logger.Log(logger.LevelError, "Error decoding program ID", logger.Error(err))
	}

	accounts, err := i0.ResolveInstructionAccounts(&transaction.Message)
	if err != nil {
		logger.Log(logger.LevelError, "Error decoding Instruction Accounts", logger.Error(err))
	}

	decodedInstruction, err := solana.DecodeInstruction(
		progKey,
		accounts,
		i0.Data,
	)
	if err != nil {
		logger.Log(logger.LevelError, "Error decoding Instructions", logger.Error(err))
	}

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
