package sell

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	lookuptable "pump_fun/app/lookup_table"
	"pump_fun/internal/core/constants"
	positionmodel "pump_fun/internal/core/position"
	"pump_fun/internal/core/tasks"
	"pump_fun/internal/monitoring/decoder"
	"pump_fun/internal/services/position"
	subscriptionhub "pump_fun/internal/services/subscription_hub"
	"pump_fun/internal/solana/client"
	"pump_fun/internal/solana/instructions"
	pumpInstructions "pump_fun/internal/solana/programs/pumpfun/instructions"
	wallets "pump_fun/internal/solana/wallet"
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
)

type Transaction struct {
	Task         *tasks.SellTask
	instructions []solana.Instruction
	transaction  *solana.Transaction
	signature    solana.Signature
}

func (st *Transaction) BuildInstructionsWithPosition(ctx context.Context, reporter subscriptionhub.TaskReporter, ps *position.Service) error {

	sellInstructions, err := getAllInstructionsForSell(st.Task, ctx, ps)
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

// func (st *Transaction) BuildInstructions(ctx context.Context, reporter subscriptionhub.TaskReporter) error {
// 	sellInstructions, err := getAllInstructionsForSell(st.Task, ctx)
// 	if err != nil {
// 		logger.Error("Error getting instructions for sell task", err)
// 		return err
// 	}

// 	if len(sellInstructions) == 0 {
// 		return fmt.Errorf("sell instruction's weren't generated properly - check sell instruction builder")
// 	}

// 	st.instructions = sellInstructions
// 	reporter.Report("Instructions Built")
// 	return nil
// }

func (st *Transaction) BuildTransaction(ctx context.Context, reporter subscriptionhub.TaskReporter) error {
	if st.Task == nil {
		return fmt.Errorf("sell task is null - check if sell task was set")
	}

	latestHash, err := client.GetLatestBlockhash(ctx)
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

func (st *Transaction) SendTransaction(ctx context.Context, reporter subscriptionhub.TaskReporter) error {
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
	txResp, err := rpcClient.SendTransaction(ctx, st.transaction)
	if err != nil {
		logger.Error(err)
		return err
	}

	st.signature = txResp
	reporter.Report(fmt.Sprintf("Tx Sent: %s", txResp))
	return nil
}

func (st *Transaction) ConfirmTransaction(ctx context.Context, reporter subscriptionhub.TaskReporter) error {
	stream := make(chan client.ConfirmMessage, 100)

	go func(stream chan client.ConfirmMessage) {
		defer close(stream)
		client.ConfirmTransactionWithStream(st.signature, ctx, stream)
	}(stream)

	for msg := range stream {
		if msg.Err != "" {
			return fmt.Errorf("%v", msg.Err)
		}
		reporter.Report(msg.Message)
	}

	return nil
}

func (st *Transaction) GetSignature() solana.Signature {
	return st.signature
}

func (st *Transaction) ExtractTokenAndSolFromTx(signature solana.Signature, ctx context.Context) (tokenAmount float64, solAmount float64, err error) {
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
		instructionData, err := base58.Decode(instruction.Data.String())
		if err != nil {
			return tokenAmount, solAmount, err
		}
		if len(instructionData) < 8 {
			continue
		}

		if !bytes.HasPrefix(instructionData, constants.SellInstructionDiscriminator[:]) {
			continue
		}

		tokenAmountInt, err := decoder.ExtractTokenAmountFromPfInstruction(instructionData)
		if err != nil {
			return tokenAmount, solAmount, err
		}

		tokenAmount = float64(tokenAmountInt)
	}

	//extract sol amount
	walletPubkey := st.Task.Wallet.PublicKey()
	var walletIndex int = -1

	for i, account := range transactionMessage.AccountKeys {
		if account.PublicKey == walletPubkey {
			walletIndex = i
		}
	}

	if walletIndex == -1 {
		return tokenAmount, solAmount, fmt.Errorf("could not find user's wallet in account keys")
	}

	preBalance := int64(tx.Meta.PostBalances[walletIndex])
	postBalance := int64(tx.Meta.PreBalances[walletIndex])

	solAmountLamport := postBalance - preBalance
	solAmount = float64(solAmountLamport)

	return tokenAmount, solAmount, nil
}

func (st *Transaction) GetTask() tasks.Task {
	return st.Task
}

func getAllInstructionsForSell(sellTask *tasks.SellTask, ctx context.Context, ps *position.Service) ([]solana.Instruction, error) {
	// we need to check if the user passed a positionId -> if so position == nil
	// else we retrieve the positonId
	var pos *positionmodel.Position
	if sellTask.Position_id != nil {
		position, err := ps.GetById(*sellTask.Position_id)
		if err != nil {
			return nil, err
		}

		isTokensRemaining := position.TokenRemaining.Cmp(big.NewFloat(0)) == 1
		if !isTokensRemaining {
			return nil, fmt.Errorf("no tokens remaining")
		}

		pos = position
	} else {
		pos = nil
	}

	computeLimitInstruction := instructions.GetComputeUnitLimitInstruction(sellTask.ComputeUnits)
	computeLimitBudgetInstruction := instructions.GetComputeUnitBudgetInstruction(sellTask.Fee, sellTask.ComputeUnits)

	sellInstructions, err := pumpInstructions.GetSellInstruction(sellTask, ctx, pos)
	if err != nil {
		logger.Error("Error getting sell instruction", err)
		return nil, err
	}

	solInstructions := []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, sellInstructions}
	return solInstructions, nil
}
