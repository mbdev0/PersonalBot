package buy

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

//TODO: instead of an interface we'll need a task_executor struct -> is this too many structs??
type TransactionSender interface {
	Execute(task *tasks.Task)
}

func Execute(buyTask *tasks.BuyTask) {
	go SendBuyTransaction(buyTask)
}

func SendBuyTransaction(buyTask *tasks.BuyTask) {
	buyTask.TransitionToNextState(tasks.TaskStateRunning)
	sendBuyTransaction(buyTask)
}

func sendBuyTransaction(buyTask *tasks.BuyTask) {
	instructions := getAllInstructionsForBuy(buyTask)
	if instructions == nil {
		logger.Error("Error creating buy instructions - no instructions created")
		return
	}

	latestHash, err := rpcclient.GetLatestBlockhash()
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
		return
	}

	accountLookupMap := lookuptable.GetAddressLookupTable()

	tx, err := solana.NewTransaction(instructions,
		latestHash.Value.Blockhash,
		solana.TransactionPayer(buyTask.Wallet.PublicKey()),
		solana.TransactionAddressTables(accountLookupMap))
	tx.Message.SetVersion(solana.MessageVersionV0)

	if err != nil {
		logger.Error("Error creating transaction", err)
		return
	}

	handlers.SignTx(tx, buyTask.Wallet)

	client := rpcclient.GetClient()
	// SIMULATE TRANSACTION
	// txResp, err := client.SimulateTransaction(context.Background(), tx)
	// if err != nil {
	// 	logger.Error("Transaction simulation failed", err)
	// 	return
	// }
	// fmt.Println(txResp.Value)

	// SEND TRANSACTION WITH OPTIONS
	// maxRetries := uint(5)
	// txResp, err := client.SendTransactionWithOpts(context.Background(), tx, rpc.TransactionOpts{Encoding: solana.EncodingBase64, SkipPreflight: false, MaxRetries: &maxRetries})
	// if err != nil {
	// 	logger.Error(err)
	// }
	// fmt.Println(txResp.String())

	// SEND TRANSACTION WITH NO OPTS
	txResp, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		logger.Error(err)
	}
	fmt.Println(txResp.String())

	buyTask.TransitionToNextState(tasks.TaskStateTransactionSent)

	isSuccess, err := rpcclient.ConfirmTransaction(txResp)
	if err != nil {
		logger.Error("Transaction confirmation failed", err)
		buyTask.TransitionToNextState(tasks.TaskStateTransactionFailed)
		return
	}
	if isSuccess {
		logger.Information("Transaction confirmed successfully")
		buyTask.TransitionToNextState(tasks.TaskStateCompleted)
	} else {
		logger.Error("Transaction confirmation failed")
		buyTask.TransitionToNextState(tasks.TaskStateTransactionFailed)
	}

}

func getAllInstructionsForBuy(buyTask *tasks.BuyTask) (buyInstructions []solana.Instruction) {

	computeLimitInstruction := instructionbuilder.GetComputeUnitLimitInstruction(buyTask.ComputeUnits)
	computeLimitBudgetInstruction := instructionbuilder.GetComputeUnitBudgetInstruction(buyTask.BuyFee, buyTask.ComputeUnits)
	idEmponenetInstruction := instructionbuilder.GetIdEmponentInstruction(buyTask.Wallet.PublicKey(), buyTask.TokenAddress)
	buyInstruction, err := instructionbuilder.GetBuyInstruction(buyTask)

	if err != nil {
		logger.Error("Error creating buy instruction", err)
		return nil
	}

	if idEmponenetInstruction != nil {
		return []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, idEmponenetInstruction, buyInstruction}
	} else {
		return []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, buyInstruction}
	}
}
