package sell

import (
	"context"
	"fmt"
	"pump_fun/internal/handlers"
	instructionbuilder "pump_fun/internal/instruction_builder"
	"pump_fun/internal/logger"
	"pump_fun/internal/models/tasks"
	rpcclient "pump_fun/internal/rpc_client"

	"github.com/gagliardetto/solana-go"
)

func SendSellTransaction(sellTask *tasks.SellTask) {
	sendSellTransaction(sellTask)
}

func sendSellTransaction(sellTask *tasks.SellTask) {
	instructions, err := getAllInstructionsForSell(sellTask)
	if err != nil {
		logger.Error("Error getting instructions for sell task", err)
		return
	}

	latestHash, err := rpcclient.GetLatestBlockhash()
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
		return
	}

	tx, err := solana.NewTransaction(instructions, latestHash.Value.Blockhash, solana.TransactionPayer(sellTask.Wallet.PublicKey()))
	if err != nil {
		logger.Error("Error creating transaction", err)
		return
	}

	handlers.SignTx(tx, sellTask.Wallet)

	client := rpcclient.GetClient()

	// simulate the transaction
	txResp, err := client.SimulateTransaction(context.Background(), tx)
	if err != nil {
		logger.Error("Transaction simulation failed", err)
		return
	}
	fmt.Println(txResp.Value)

	// SEND TRANSACTION WITH OPTIONS
	// maxRetries := uint(5)
	// txResp, err := client.SendTransactionWithOpts(context.Background(), tx, rpc.TransactionOpts{Encoding: solana.EncodingBase64, SkipPreflight: true, MaxRetries: &maxRetries})
	// if err != nil {
	// 	logger.Error("Error sending transaction", err)
	// 	return
	// }
	// fmt.Println(txResp.String())

	// SEND TRANSACTION WITH NO OPTS
	// txResp, err := client.SendTransaction(context.Background(), tx)
	// if err != nil {
	// 	logger.Error(err)
	// }
	// fmt.Println(txResp.String())
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
