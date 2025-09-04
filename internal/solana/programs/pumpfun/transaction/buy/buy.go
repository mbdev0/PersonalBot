package buy

import (
	"bytes"
	"context"
	"fmt"
	lookuptable "pump_fun/app/lookup_table"
	"pump_fun/internal/core/constants"
	"pump_fun/internal/core/tasks"
	"pump_fun/internal/monitoring/decoder"
	subscriptionhub "pump_fun/internal/services/subscription_hub"
	"pump_fun/internal/solana/client"
	"pump_fun/internal/solana/instructions"
	pumpInstructions "pump_fun/internal/solana/programs/pumpfun/instructions"
	wallets "pump_fun/internal/solana/wallet"
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58/base58"
)

type Transaction struct {
	BuyTask      *tasks.BuyTask
	instructions *[]solana.Instruction
	transaction  *solana.Transaction
	signature    solana.Signature
}

func (bt *Transaction) BuildInstructions(ctx context.Context, reporter subscriptionhub.TaskReporter) error {
	if bt.BuyTask == nil {
		return fmt.Errorf("buy task was nil - make sure buy task is set")
	}

	buyInstructions, err := bt.getAllInstructionsForBuy(bt.BuyTask, ctx)
	if buyInstructions == nil || err != nil {
		logger.Error("Error creating buy instructions - no instructions created")
		return err
	}

	bt.instructions = &buyInstructions
	reporter.Report("Instructions Built")
	return nil
}

func (bt *Transaction) BuildTransaction(ctx context.Context, reporter subscriptionhub.TaskReporter) error {
	if bt.BuyTask == nil {
		return fmt.Errorf("buy task was nil - make sure buy task is set")
	}

	latestHash, err := client.GetLatestBlockhash(ctx)
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

func (bt *Transaction) SendTransaction(ctx context.Context, reporter subscriptionhub.TaskReporter) error {
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
	txResp, err := rpcClient.SendTransactionWithOpts(ctx, bt.transaction, rpc.TransactionOpts{Encoding: solana.EncodingBase64, SkipPreflight: true})
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

func (bt *Transaction) ConfirmTransaction(ctx context.Context, reporter subscriptionhub.TaskReporter) error {
	stream := make(chan client.ConfirmMessage, 100)

	go func(stream chan client.ConfirmMessage) {
		defer close(stream)
		client.ConfirmTransactionWithStream(bt.signature, ctx, stream)
	}(stream)

	for msg := range stream {
		if msg.Err != "" {
			return fmt.Errorf("%v", msg.Err)
		}
		reporter.Report(msg.Message)
	}

	tokenAmnt, solAmnt, err := bt.extractTokenAndSolFromTx(bt.signature, ctx)
	if err != nil {
		return err
	}
	fmt.Println(tokenAmnt, solAmnt)
	return nil
}

func (bt *Transaction) GetTask() tasks.Task {
	return bt.BuyTask
}

func (bt *Transaction) extractTokenAndSolFromTx(signature solana.Signature, ctx context.Context) (tokenAmount float64, solAmount float64, err error) {
	solClient := client.GetClient()
	tx, err := solClient.GetParsedTransaction(ctx, signature, &rpc.GetParsedTransactionOpts{Commitment: rpc.CommitmentConfirmed, MaxSupportedTransactionVersion: &rpc.MaxSupportedTransactionVersion0})
	if err != nil {
		return tokenAmount, solAmount, err
	}

	if tx.Meta.Err != nil {
		return tokenAmount, solAmount, fmt.Errorf("error in transaction whilst extracting token amount + sol amount")
	}

	transactionMessage := tx.Transaction.Message

	instructions := transactionMessage.Instructions
	//extract token amount
	for _, instruction := range instructions {
		logger.Information(base58.Decode(instruction.Data.String()))
		instructionData, err := base58.Decode(instruction.Data.String())
		if err != nil {
			return tokenAmount, solAmount, err
		}
		if len(instructionData) < 8 {
			continue
		}

		if !bytes.HasPrefix(instructionData, constants.BuyInstructionDiscriminator[:]) {
			continue
		}

		tokenAmountInt, err := decoder.ExtractTokenAmountFromBuyInstruction(instructionData)
		if err != nil {
			return tokenAmount, solAmount, err
		}

		tokenAmount = float64(tokenAmountInt) / constants.TokenAmountDecimals
	}

	//extract sol amount
	walletPubkey := bt.BuyTask.Wallet.PublicKey()
	var walletIndex int = -1

	for i, account := range transactionMessage.AccountKeys {
		if account.PublicKey == walletPubkey {
			walletIndex = i
		}
	}

	if walletIndex == -1 {
		return tokenAmount, solAmount, fmt.Errorf("could not find user's wallet in account keys")
	}

	solAmountLamport := tx.Meta.PreBalances[walletIndex] - tx.Meta.PostBalances[walletIndex]
	solAmount = float64(solAmountLamport) / float64(solana.LAMPORTS_PER_SOL)

	return tokenAmount, solAmount, nil
}

func (bt *Transaction) getAllInstructionsForBuy(buyTask *tasks.BuyTask, ctx context.Context) (buyInstructions []solana.Instruction, err error) {

	computeLimitInstruction := instructions.GetComputeUnitLimitInstruction(buyTask.ComputeUnits)
	computeLimitBudgetInstruction := instructions.GetComputeUnitBudgetInstruction(buyTask.Fee, buyTask.ComputeUnits)
	idEmponenetInstruction, err := instructions.GetIdempotentInstruction(buyTask.Wallet.PublicKey(), buyTask.Token, ctx)
	if err != nil {
		return nil, err
	}

	buyInstruction, err := pumpInstructions.GetBuyInstruction(buyTask, ctx)

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
