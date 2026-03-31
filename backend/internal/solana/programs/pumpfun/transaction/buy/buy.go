package buy

import (
	"bytes"
	"context"
	"fmt"
	lookuptable "personal_bot/app/lookup_table"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/monitoring/decoder"
	"personal_bot/internal/services/position"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	"personal_bot/internal/solana/client"
	"personal_bot/internal/solana/instructions"
	pumpInstructions "personal_bot/internal/solana/programs/pumpfun/instructions"
	wallets "personal_bot/internal/solana/wallet"
	"personal_bot/pkg/logger"

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

func (bt *Transaction) BuildInstructionsWithPosition(ctx context.Context, publisher subscriptionhub.Publisher, ps *position.Service) error {
	return bt.buildInstructions(ctx, publisher)
}

func (bt *Transaction) buildInstructions(ctx context.Context, publisher subscriptionhub.Publisher) error {
	if bt.BuyTask == nil {
		return fmt.Errorf("buy task was nil - make sure buy task is set")
	}

	buyInstructions, err := bt.getAllInstructionsForBuy(bt.BuyTask, ctx)
	if buyInstructions == nil || err != nil {
		logger.Error("Error creating buy instructions - no instructions created")
		return err
	}

	bt.instructions = &buyInstructions
	publisher.PublishMessage(bt.BuyTask, "Instructions Built")

	return nil
}

func (bt *Transaction) BuildTransaction(ctx context.Context, publisher subscriptionhub.Publisher) error {
	if bt.BuyTask == nil {
		return fmt.Errorf("buy task was nil - make sure buy task is set")
	}

	latestHash, err := client.GetLatestBlockhash(bt.GetTask().HttpClient(), ctx)
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
		return err
	}

	opts := []solana.TransactionOption{
		solana.TransactionPayer(bt.BuyTask.Wallet.PublicKey()),
	}

	accountLookupMap, err := lookuptable.GetAddressLookupTable(bt.GetTask().HttpClient())
	if err != nil {
		logger.Error("Error getting address lookup table, proceeding without it: ", err)
	} else {
		opts = append(opts, solana.TransactionAddressTables(accountLookupMap))
	}

	tx, err := solana.NewTransaction(*bt.instructions, latestHash.Value.Blockhash, opts...)
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
	publisher.PublishMessage(bt.BuyTask, "TX Built")

	return nil
}

func (bt *Transaction) SendTransaction(ctx context.Context, publisher subscriptionhub.Publisher) error {
	rpcClient := bt.GetTask().HttpClient()
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
	publisher.PublishMessage(bt.BuyTask, fmt.Sprintf("Tx Sent: %s", txResp))
	return nil
}

func (bt *Transaction) ConfirmTransaction(ctx context.Context, publisher subscriptionhub.Publisher) error {
	rpcClient := bt.GetTask().HttpClient()

	stream := make(chan client.ConfirmMessage, 100)

	go func(stream chan client.ConfirmMessage) {
		defer close(stream)
		client.ConfirmTransactionWithStream(rpcClient, bt.signature, ctx, stream)
	}(stream)

	for msg := range stream {
		if msg.Err != "" {
			publisher.PublishMessage(bt.BuyTask, msg.Err)
			return fmt.Errorf("%v", msg.Err)
		}
		publisher.PublishMessage(bt.BuyTask, msg.Message)
	}

	tokenAmnt, solAmnt, err := bt.ExtractTokenAndSolFromTx(bt.signature, ctx)
	if err != nil {
		return err
	}

	fmt.Println(tokenAmnt, solAmnt)

	return nil
}

func (bt *Transaction) GetTask() tasks.Task {
	return bt.BuyTask
}

func (bt *Transaction) GetSignature() solana.Signature {
	return bt.signature
}

func (bt *Transaction) ExtractTokenAndSolFromTx(signature solana.Signature, ctx context.Context) (tokenAmount float64, solAmount float64, err error) {
	solClient := bt.GetTask().HttpClient()
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
			if len(instructionData) == 0 {
				continue
			}
			return tokenAmount, solAmount, err
		}
		if len(instructionData) < 8 {
			continue
		}

		if !bytes.HasPrefix(instructionData, constants.BuyInstructionDiscriminator[:]) {
			continue
		}

		tokenAmountInt, err := decoder.ExtractTokenAmountFromPfInstruction(instructionData)
		if err != nil {
			return tokenAmount, solAmount, err
		}

		tokenAmount = float64(tokenAmountInt)
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
	solAmount = float64(solAmountLamport)

	return tokenAmount, solAmount, nil
}

func (bt *Transaction) getAllInstructionsForBuy(buyTask *tasks.BuyTask, ctx context.Context) (buyInstructions []solana.Instruction, err error) {

	computeLimitInstruction := instructions.GetComputeUnitLimitInstruction(buyTask.ComputeUnits)
	computeLimitBudgetInstruction := instructions.GetComputeUnitBudgetInstruction(buyTask.Fee, buyTask.ComputeUnits)
	idEmponenetInstruction, err := instructions.GetIdempotentInstruction(buyTask.Wallet.PublicKey(), buyTask.Token, ctx, buyTask.HttpClient())
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
