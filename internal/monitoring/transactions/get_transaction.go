package transactions

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	solanaclient "pump_fun/internal/solana-client"

	bin "github.com/gagliardetto/binary"

	"pump_fun/internal/models"
)

func GetTransaction(signature string) interface{}{
	programID := solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")
    solana.RegisterInstructionDecoder(programID, CustomInstructionDecoder)
	transaction := getTransaction(signature)
	return decodeCreateTransaction(transaction)
}

func getTransaction(signature string) (*solana.Transaction){
	rpcClient := solanaclient.NewHttpClient()

	version := uint64(0)
	out, err := rpcClient.GetTransaction(
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
		logger.Log(slog.LevelError, "Error getting transaction", slog.String("error: ", err.Error()))
	}

	transaction, err := solana.TransactionFromDecoder(bin.NewBinDecoder(out.Transaction.GetBinary()))
    if err != nil {
		logger.Log(slog.LevelError, "Error decoding transaction", slog.String("error: ", err.Error()))
    }
	return transaction
}

func decodeCreateTransaction(transaction *solana.Transaction) models.DecodedInstruction {
	i0 := transaction.Message.Instructions[3]
	progKey, err := transaction.ResolveProgramIDIndex(i0.ProgramIDIndex)
	if err != nil {
		logger.Log(slog.LevelError, "Error decoding program ID", slog.String("error: ", err.Error()))
	}

	accounts, err := i0.ResolveInstructionAccounts(&transaction.Message)
	if err != nil {
		logger.Log(slog.LevelError, "Error decoding Instruction Accounts", slog.String("error: ", err.Error()))
	}
  
	decodedInstruction, err := solana.DecodeInstruction(
		progKey,
		accounts,
		i0.Data,
	  )
	  if err != nil {
		logger.Log(slog.LevelError, "Error decoding Instructions", slog.String("error: ", err.Error()))
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
        Name:   decodedInstructionMap["Name"],
        Symbol: decodedInstructionMap["Symbol"],
        IPFS_URL:    decodedInstructionMap["Uri"],
    }
}
