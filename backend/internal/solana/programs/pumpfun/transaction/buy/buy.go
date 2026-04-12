package buy

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	lookuptable "personal_bot/app/lookup_table"
	"personal_bot/internal/core/constants"
	positionmodels "personal_bot/internal/core/position"
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
	BuyTask         *tasks.BuyTask
	instructions    *[]solana.Instruction
	transaction     *solana.Transaction
	signature       solana.Signature
	PositionService *position.Service
}

func (bt *Transaction) BuildInstructions(ctx context.Context, publisher subscriptionhub.Publisher) error {
	if bt.BuyTask == nil {
		return fmt.Errorf("buy task was nil - make sure buy task is set")
	}

	buyInstructions, err := bt.getAllInstructionsForBuy(ctx, bt.BuyTask)
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

	latestHash, err := client.GetLatestBlockhash(ctx, bt.getHttpClient())
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
		return err
	}

	opts := []solana.TransactionOption{
		solana.TransactionPayer(bt.BuyTask.Wallet.PublicKey()),
	}

	accountLookupMap, err := lookuptable.GetAddressLookupTable(bt.getHttpClient())
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
	rpcClient := bt.getHttpClient()
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
	rpcClient := bt.getHttpClient()

	stream := make(chan client.ConfirmMessage, 100)

	go func(stream chan client.ConfirmMessage) {
		defer close(stream)
		client.ConfirmTransactionWithStream(ctx, rpcClient, bt.signature, stream)
	}(stream)

	for msg := range stream {
		if msg.Err != "" {
			publisher.PublishMessage(bt.BuyTask, msg.Err)
			return fmt.Errorf("%v", msg.Err)
		}
		publisher.PublishMessage(bt.BuyTask, msg.Message)
	}

	return nil
}

func (bt *Transaction) GetTask() tasks.Task {
	return bt.BuyTask
}

func (bt *Transaction) getHttpClient() *rpc.Client {
	return bt.BuyTask.HttpClient()
}

func (bt *Transaction) UpdatePosition(ctx context.Context, publisher subscriptionhub.Publisher) (tokenAmount, solAmount float64, pos *positionmodels.Position, err error) {
	tokenAmnt, solAmnt, err := bt.extractTokenAndSolFromTx(ctx, bt.signature)
	if err != nil {
		return
	}

	err = bt.PositionService.ReportBuy(ctx, bt.BuyTask.Id(), bt.BuyTask.Token, bt.BuyTask.Wallet.PublicKey(), new(big.Float).SetFloat64(tokenAmnt), new(big.Float).SetFloat64(solAmnt))
	if err != nil {
		return
	}

	return tokenAmnt, solAmnt, nil, nil
}

func (bt *Transaction) extractTokenAndSolFromTx(ctx context.Context, signature solana.Signature) (tokenAmount float64, solAmount float64, err error) {
	solClient := bt.getHttpClient()
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
	walletIndex := -1

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

func (bt *Transaction) GetSignature() string {
	return bt.signature.String()
}

func (bt *Transaction) getAllInstructionsForBuy(ctx context.Context, buyTask *tasks.BuyTask) (buyInstructions []solana.Instruction, err error) {

	computeLimitInstruction := instructions.GetComputeUnitLimitInstruction(buyTask.ComputeUnits)
	computeLimitBudgetInstruction := instructions.GetComputeUnitBudgetInstruction(buyTask.Fee, buyTask.ComputeUnits)
	idEmponenetInstruction, err := instructions.GetIdempotentInstruction(ctx, buyTask.Wallet.PublicKey(), buyTask.Token, buyTask.HttpClient())
	if err != nil {
		return nil, err
	}

	buyInstruction, err := pumpInstructions.GetBuyInstruction(ctx, buyTask)

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
