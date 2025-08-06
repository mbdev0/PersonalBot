package sell

import (
	"fmt"
	"pump_fun/internal/handlers"
	instructionbuilder "pump_fun/internal/instruction_builder"
	lookuptable "pump_fun/internal/launch/lookup_table"
	"pump_fun/internal/models/tasks"
	rpcclient "pump_fun/internal/rpc_client"
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go"
)

type SellTransaction struct {
	Task         *tasks.SellTask
	instructions []solana.Instruction
	transaction  *solana.Transaction
	signature    solana.Signature
}

func (st *SellTransaction) BuildInstructions() error {
	sellInstructions, err := getAllInstructionsForSell(st.Task)
	if err != nil {
		logger.Error("Error getting instructions for sell task", err)
		return err
	}

	if len(sellInstructions) == 0 {
		return fmt.Errorf("sell instruction's weren't generated properly - check sell instruction builder")
	}

	st.instructions = sellInstructions
	return nil
}

func (st *SellTransaction) BuildTransaction() error {
	if st.Task == nil {
		return fmt.Errorf("sell task is null - check if sell task was set")
	}

	latestHash, err := rpcclient.GetLatestBlockhash(st.Task.CancelToken)
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
		return err
	}

	accountLookupMap := lookuptable.GetAddressLookupTable()
	tx, err := solana.NewTransaction(st.instructions,
		latestHash.Value.Blockhash,
		solana.TransactionPayer(st.Task.Wallet.PublicKey()),
		solana.TransactionAddressTables(accountLookupMap))

	if err != nil {
		logger.Error("Error creating transaction", err)
		return err
	}

	handlers.SignTx(tx, st.Task.Wallet)
	st.transaction = tx
	return nil
}

func (st *SellTransaction) SendTransaction() error {
	client := rpcclient.GetClient()

	// simulate the transaction
	// txResp, err := client.SimulateTransaction(st.Task.CancelToken.CancellationContext, st.transaction)
	// if err != nil {
	// 	logger.Error("Transaction simulation failed", err)
	// 	return
	// }
	// fmt.Println(txResp.Value)

	// SEND TRANSACTION WITH OPTIONS
	// maxRetries := uint(5)
	// txResp, err := client.SendTransactionWithOpts(st.Task.CancelToken.CancellationContext, st.transaction, rpc.TransactionOpts{Encoding: solana.EncodingBase64, SkipPreflight: true, MaxRetries: &maxRetries})
	// if err != nil {
	// 	logger.Error("Error sending transaction", err)
	// 	return
	// }
	// fmt.Println(txResp.String())

	// SEND TRANSACTION WITH NO OPTS
	txResp, err := client.SendTransaction(st.Task.CancelToken.CancellationContext, st.transaction)
	if err != nil {
		logger.Error(err)
		return err
	}

	st.signature = txResp
	return nil
}

func (st *SellTransaction) ConfirmTransaction() error {
	isSuccess, err := rpcclient.ConfirmTransaction(st.signature, st.Task.CancelToken)
	if err != nil {
		logger.Error("Transaction confirmation failed", err)
		return err
	}
	if isSuccess {
		logger.Information("Transaction confirmed successfully")
	} else {
		logger.Error("Transaction confirmation failed")
		return fmt.Errorf("transaction confirmation failed")
	}
	return nil
}

func (st *SellTransaction) GetTask() tasks.Task {
	return st.Task
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
