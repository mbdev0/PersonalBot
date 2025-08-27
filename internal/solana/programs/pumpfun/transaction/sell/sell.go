package sell

import (
	"fmt"
	lookuptable "pump_fun/app/lookup_table"
	"pump_fun/internal/core/tasks"
	subscriptionhub "pump_fun/internal/services/subscription_hub"
	"pump_fun/internal/solana/client"
	"pump_fun/internal/solana/instructions"
	pumpInstructions "pump_fun/internal/solana/programs/pumpfun/instructions"
	wallets "pump_fun/internal/solana/wallet"
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go"
)

type Transaction struct {
	Task         *tasks.SellTask
	instructions []solana.Instruction
	transaction  *solana.Transaction
	signature    solana.Signature
}

func (st *Transaction) BuildInstructions(reporter subscriptionhub.TaskReporter) error {
	sellInstructions, err := getAllInstructionsForSell(st.Task)
	if err != nil {
		logger.Error("Error getting instructions for sell task", err)
		return err
	}

	if len(sellInstructions) == 0 {
		return fmt.Errorf("sell instruction's weren't generated properly - check sell instruction builder")
	}

	st.instructions = sellInstructions
	reporter.Report("Instructions Built")
	return nil
}

func (st *Transaction) BuildTransaction(reporter subscriptionhub.TaskReporter) error {
	if st.Task == nil {
		return fmt.Errorf("sell task is null - check if sell task was set")
	}

	latestHash, err := client.GetLatestBlockhash(st.Task.Ctx())
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

	err = wallets.SignTx(tx, st.Task.Wallet)
	if err != nil {
		return err
	}

	st.transaction = tx
	reporter.Report("Tx Built")

	return nil
}

func (st *Transaction) SendTransaction(reporter subscriptionhub.TaskReporter) error {
	rpcClient := client.GetClient()

	// simulate the transaction
	// txResp, err := rpcClient.SimulateTransaction(st.Task.Ctx(), st.transaction)
	// if err != nil {
	// 	logger.Error("Transaction simulation failed", err)
	// 	return
	// }
	// fmt.Println(txResp.Value)

	// SEND TRANSACTION WITH OPTIONS
	// maxRetries := uint(5)
	// txResp, err := rpcClient.SendTransactionWithOpts(st.Task.Ctx(), st.transaction, rpc.TransactionOpts{Encoding: solana.EncodingBase64, SkipPreflight: true, MaxRetries: &maxRetries})
	// if err != nil {
	// 	logger.Error("Error sending transaction", err)
	// 	return
	// }
	// fmt.Println(txResp.String())

	// SEND TRANSACTION WITH NO OPTS
	txResp, err := rpcClient.SendTransaction(st.Task.Ctx(), st.transaction)
	if err != nil {
		logger.Error(err)
		return err
	}

	st.signature = txResp
	reporter.Report(fmt.Sprintf("Tx Sent: %s", txResp))
	return nil
}

func (st *Transaction) ConfirmTransaction(reporter subscriptionhub.TaskReporter) error {
	stream := make(chan client.ConfirmMessage, 100)

	go func(stream chan client.ConfirmMessage) {
		defer close(stream)
		client.ConfirmTransactionWithStream(st.signature, st.Task.Ctx(), stream)
	}(stream)

	for msg := range stream {
		if msg.Err != "" {
			return fmt.Errorf(msg.Err)
		}
		reporter.Report(msg.Message)
	}

	return nil
}
func (st *Transaction) GetTask() tasks.Task {
	return st.Task
}

func getAllInstructionsForSell(sellTask *tasks.SellTask) ([]solana.Instruction, error) {
	computeLimitInstruction := instructions.GetComputeUnitLimitInstruction(sellTask.ComputeUnits)
	computeLimitBudgetInstruction := instructions.GetComputeUnitBudgetInstruction(sellTask.Fee, sellTask.ComputeUnits)

	sellInstructions, err := pumpInstructions.GetSellInstruction(sellTask)
	if err != nil {
		logger.Error("Error getting sell instruction", err)
		return nil, err
	}

	solInstructions := []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, sellInstructions}
	return solInstructions, nil
}
