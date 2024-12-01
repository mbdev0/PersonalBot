package transactions

import (
	"bytes"
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"

	"pump_fun/internal/constants"
	"pump_fun/internal/models"
	"pump_fun/internal/monitoring/geyser"
	solanaclient "pump_fun/internal/solana-client"

	bin "github.com/gagliardetto/binary"

	"pump_fun/internal/logger"
)

func GetTransaction(signature string) (*models.Coin, error) {
	return nil, nil

	programID := solana.MustPublicKeyFromBase58(constants.ProgramID)
	solana.RegisterInstructionDecoder(programID, CustomInstructionDecoder)
	transaction, err := getTransaction(signature)

	if err != nil {
		return nil, err
	}

	return ParseTransaction(transaction), nil
}

func DecryptTransactionNotification(transaction geyser.TransactionNotification, coinStructChan chan models.Coin) {
	// Decrypt transaction
	// we get the entire transaction notification
	// get the 4th instruction if it exists
	// check if the first 8 bytes are equal to the create discriminator
	// if it is, skip the first 8 bytes and decode the rest of the instruction
	// return the decoded instruction
	// send the decoded instruction to the coinStructChan

	// get the 4th instruction if it exists
	if len(transaction.Params.Result.Transaction.TransactionDetails.Message.Instructions) < 4 {
		return
	}

	compiled_instruction := transaction.Params.Result.Transaction.TransactionDetails.Message.Instructions[3]
	accounts := transaction.Params.Result.Transaction.TransactionDetails.Message.AccountKeys

	if len(compiled_instruction.Data) < 8 {
		return
	}

	//decode instruction from base58
	instruction, err := base58.Decode(compiled_instruction.Data)
	if err != nil {
		logger.Log(logger.LevelError, "Error decoding instruction", logger.Error(err))
		return
	}

	discriminator := instruction[:8]

	if bytes.Equal(discriminator, constants.CreateInstructionDiscriminator[:]) {
		decodedInstruction, err := CreateInstructionDecoderNew(&accounts, instruction)
		if err != nil {
			logger.Log(logger.LevelError, "Error decoding instruction", logger.Error(err))
			return
		}
		// fmt.Println(mapToStruct(decodedInstruction))
		coinStructChan <- mapToStruct(decodedInstruction)
	}

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
func ParseTransaction(transaction *solana.Transaction) *models.Coin {
	compiledInstruction := transaction.Message.Instructions[3]
	decodedInstruction := DecodeInstruction(compiledInstruction, transaction)
	decodedInstructionStruct := mapToStruct(decodedInstruction)
	return &decodedInstructionStruct
}

func mapToStruct(decodedInstruction interface{}) models.Coin {
	decodedInstructionMap, ok := decodedInstruction.(map[string]string)
	if !ok {
		panic("decodedInstruction is not of type map[string]string")
	}

	/*
		Signature        string
			Name             string
			Symbol           string
			IPFS_URL         string
			TokenAddr        string
			CreatorAddr      string
			DevHoldingAmount float64
	*/

	return models.Coin{
		CoinData: models.MintData{
			Name:             decodedInstructionMap["Name"],
			Symbol:           decodedInstructionMap["Symbol"],
			IPFS_URL:         decodedInstructionMap["IPFS_URL"],
			TokenAddr:        "",
			CreatorAddr:      "",
			DevHoldingAmount: 0,
		},
		IPFSData: models.IPFS{
			TelegramURL: "https://t.me/pumpfun",
			TwitterURL:  "https://twitter.com",
			WebsiteURL:  "https://pump.fun",
			ImageURL:    "https://pump.fun/pumpfun.png",
		},
	}

	// return models.DecodedInstruction{
	// 	Name:     decodedInstructionMap["Name"],
	// 	Symbol:   decodedInstructionMap["Symbol"],
	// 	IPFS_URL: decodedInstructionMap["IPFS_URL"],
	// }
}
