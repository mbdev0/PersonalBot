package sell

import (
	"context"
	"fmt"
	"pump_fun/internal/handlers"
	instructionbuilder "pump_fun/internal/instruction_builder"
	lookuptable "pump_fun/internal/launch/lookup_table"
	"pump_fun/internal/models/tasks"
	rpcclient "pump_fun/internal/rpc_client"
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go"
)

func SendSellTransaction(sellTask *tasks.SellTask) (signature string, err error) {
	signature, err = sendSellTransaction(sellTask)
	if err != nil {
		return "", err
	}

	return signature, nil
}

func sendSellTransaction(sellTask *tasks.SellTask) (signature string, err error) {
	instructions, err := getAllInstructionsForSell(sellTask)
	if err != nil {
		logger.Error("Error getting instructions for sell task", err)
		return "", err
	}

	latestHash, err := rpcclient.GetLatestBlockhash()
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
		return "", err
	}

	accountLookupMap := lookuptable.GetAddressLookupTable()

	tx, err := solana.NewTransaction(instructions,
		latestHash.Value.Blockhash,
		solana.TransactionPayer(sellTask.Wallet.PublicKey()),
		solana.TransactionAddressTables(accountLookupMap))

	if err != nil {
		logger.Error("Error creating transaction", err)
		return "", err
	}

	handlers.SignTx(tx, sellTask.Wallet)

	client := rpcclient.GetClient()

	// simulate the transaction
	// txResp, err := client.SimulateTransaction(context.Background(), tx)
	// if err != nil {
	// 	logger.Error("Transaction simulation failed", err)
	// 	return
	// }
	// fmt.Println(txResp.Value)

	// SEND TRANSACTION WITH OPTIONS
	// maxRetries := uint(5)
	// txResp, err := client.SendTransactionWithOpts(context.Background(), tx, rpc.TransactionOpts{Encoding: solana.EncodingBase64, SkipPreflight: true, MaxRetries: &maxRetries})
	// if err != nil {
	// 	logger.Error("Error sending transaction", err)
	// 	return
	// }
	// fmt.Println(txResp.String())

	// SEND TRANSACTION WITH NO OPTS
	txResp, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	fmt.Println(txResp.String())

	return txResp.String(), err
}

func getAllInstructionsForSell(sellTask *tasks.SellTask) ([]solana.Instruction, error) {
	computeLimitInstruction := instructionbuilder.GetComputeUnitLimitInstruction(sellTask.ComputeUnits)
	computeLimitBudgetInstruction := instructionbuilder.GetComputeUnitBudgetInstruction(sellTask.SellFee, sellTask.ComputeUnits)

	sellInstructions, err := instructionbuilder.GetSellInstruction(sellTask)
	if err != nil {
		logger.Error("Error getting sell instruction", err)
		return nil, err
	}

	instructions := []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, sellInstructions}
	return instructions, nil
}
