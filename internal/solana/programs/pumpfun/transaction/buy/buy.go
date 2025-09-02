package buy

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
	"github.com/gagliardetto/solana-go/rpc"
)

type Transaction struct {
	BuyTask      *tasks.BuyTask
	instructions *[]solana.Instruction
	transaction  *solana.Transaction
	signature    solana.Signature
}

func (bt *Transaction) BuildInstructions(reporter subscriptionhub.TaskReporter) error {
	if bt.BuyTask == nil {
		return fmt.Errorf("buy task was nil - make sure buy task is set")
	}

	buyInstructions, err := getAllInstructionsForBuy(bt.BuyTask)
	if buyInstructions == nil || err != nil {
		logger.Error("Error creating buy instructions - no instructions created")
		return err
	}

	bt.instructions = &buyInstructions
	reporter.Report("Instructions Built")
	return nil
}

func (bt *Transaction) BuildTransaction(reporter subscriptionhub.TaskReporter) error {
	if bt.BuyTask == nil {
		return fmt.Errorf("buy task was nil - make sure buy task is set")
	}

	latestHash, err := client.GetLatestBlockhash(bt.BuyTask.Ctx())
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
		return err
	}

	accountLookupMap := lookuptable.GetAddressLookupTable()
	tx, err := solana.NewTransaction(*bt.instructions,
		latestHash.Value.Blockhash,
		solana.TransactionPayer(bt.BuyTask.Wallet.PublicKey()),
		solana.TransactionAddressTables(accountLookupMap))

	if err != nil {
		logger.Error("Error creating transaction", err)
		return err
	}

	tx.Message.SetVersion(solana.MessageVersionV0)

	err = wallets.SignTx(tx, bt.BuyTask.Wallet)
	if err != nil {
		return err
	}
	bt.transaction = tx
	reporter.Report("Tx Built")

	return nil
}

func (bt *Transaction) SendTransaction(reporter subscriptionhub.TaskReporter) error {
	rpcClient := client.GetClient()
	// SIMULATE TRANSACTION
	// txResp, err := rpcClient.SimulateTransaction(bt.BuyTask.Ctx(), bt.transaction)
	// if err != nil {
	// 	logger.Error("Transaction simulation failed", err)
	// 	return nil
	// }
	// fmt.Println(txResp.Value)

	// SEND TRANSACTION WITH OPTIONS
	// maxRetries := uint(5)
	txResp, err := rpcClient.SendTransactionWithOpts(bt.BuyTask.Ctx(), bt.transaction, rpc.TransactionOpts{Encoding: solana.EncodingBase64, SkipPreflight: true})
	if err != nil {
		logger.Error(err)
		return err
	}
	fmt.Println(txResp.String())

	// SEND TRANSACTION WITH NO OPTS
	// fmt.Println(bt.transaction.String())
	// txResp, err := rpcClient.SendTransaction(bt.BuyTask.Ctx(), bt.transaction)
	// if err != nil {
	// 	logger.Error(err)
	// 	return err
	// }

	bt.signature = txResp
	reporter.Report(fmt.Sprintf("Tx Sent: %s", txResp))
	return nil
}

func (bt *Transaction) ConfirmTransaction(reporter subscriptionhub.TaskReporter) error {
	stream := make(chan client.ConfirmMessage, 100)

	go func(stream chan client.ConfirmMessage) {
		defer close(stream)
		client.ConfirmTransactionWithStream(bt.signature, bt.BuyTask.Ctx(), stream)
	}(stream)

	for msg := range stream {
		if msg.Err != "" {
			return fmt.Errorf(msg.Err)
		}
		reporter.Report(msg.Message)
	}

	return nil
}

func (bt *Transaction) GetTask() tasks.Task {
	return bt.BuyTask
}

func getAllInstructionsForBuy(buyTask *tasks.BuyTask) (buyInstructions []solana.Instruction, err error) {

	computeLimitInstruction := instructions.GetComputeUnitLimitInstruction(buyTask.ComputeUnits)
	computeLimitBudgetInstruction := instructions.GetComputeUnitBudgetInstruction(buyTask.Fee, buyTask.ComputeUnits)
	idEmponenetInstruction, err := instructions.GetIdempotentInstruction(buyTask.Wallet.PublicKey(), buyTask.Token, buyTask.Ctx())
	if err != nil {
		return nil, err
	}

	buyInstruction, err := pumpInstructions.GetBuyInstruction(buyTask)

	if err != nil {
		logger.Error("Error creating buy instruction", err)
		return nil, err
	}

	if idEmponenetInstruction != nil {
		return []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, idEmponenetInstruction, buyInstruction}, nil
	} else {
		return []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, buyInstruction}, nil
	}
}
