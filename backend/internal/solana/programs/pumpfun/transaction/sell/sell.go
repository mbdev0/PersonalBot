package sell

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	lookuptable "personal_bot/app/lookup_table"
	"personal_bot/internal/core/constants"
	positionmodel "personal_bot/internal/core/position"
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
	"github.com/mr-tron/base58"
)

type Transaction struct {
	Task            *tasks.SellTask
	instructions    []solana.Instruction
	transaction     *solana.Transaction
	signature       solana.Signature
	PositionService *position.Service
}

func (st *Transaction) BuildInstructions(ctx context.Context, publisher subscriptionhub.Publisher) error {

	sellInstructions, err := st.getAllInstructionsForSell(ctx, st.Task)
	if err != nil {
		logger.Error("Error getting instructions for sell task", err)
		return err
	}

	if len(sellInstructions) == 0 {
		return fmt.Errorf("sell instruction's weren't generated properly - check sell instruction builder")
	}

	st.instructions = sellInstructions
	publisher.PublishMessage(st.Task, "instructions built")
	return nil
}

func (st *Transaction) BuildTransaction(ctx context.Context, publisher subscriptionhub.Publisher) error {
	if st.Task == nil {
		return fmt.Errorf("sell task is null - check if sell task was set")
	}

	latestHash, err := client.GetLatestBlockhash(ctx, st.getHttpClient())
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
		return err
	}

	opts := []solana.TransactionOption{
		solana.TransactionPayer(st.Task.Wallet.PublicKey()),
	}

	accountLookupMap, err := lookuptable.GetAddressLookupTable(st.getHttpClient())
	if err != nil {
		logger.Error("Error getting address lookup table, proceeding without it: ", err)
	} else {
		opts = append(opts, solana.TransactionAddressTables(accountLookupMap))
	}

	tx, err := solana.NewTransaction(st.instructions, latestHash.Value.Blockhash, opts...)
	if err != nil {
		logger.Error("Error creating transaction", err)
		return err
	}

	err = wallets.SignTx(tx, st.Task.Wallet)
	if err != nil {
		return err
	}

	st.transaction = tx
	publisher.PublishMessage(st.Task, "tx built")

	return nil
}

func (st *Transaction) SendTransaction(ctx context.Context, publisher subscriptionhub.Publisher) error {
	rpcClient := st.getHttpClient()

	// simulate the transaction
	// txResp, err := rpcClient.SimulateTransaction(st.Task.Ctx(), st.transaction)
	// if err != nil {
	// 	logger.Error("Transaction simulation failed", err)
	// 	return
	// }
	// fmt.Println(txResp.Value)

	// SEND TRANSACTION WITH OPTIONS
	// maxRetries := uint(5)
	txResp, err := rpcClient.SendTransactionWithOpts(ctx, st.transaction, rpc.TransactionOpts{Encoding: solana.EncodingBase64, SkipPreflight: true})
	if err != nil {
		logger.Error("Error sending transaction", err)
		publisher.PublishMessage(st.Task, "failure whilst sending transaction")
		return err
	}
	fmt.Println(txResp.String())

	// SEND TRANSACTION WITH NO OPTS
	// txResp, err := rpcClient.SendTransaction(ctx, st.transaction)
	// if err != nil {
	// 	logger.Error(err)
	// 	publisher.PublishMessage(st.Task, "failure whilst sending transaction")
	// 	return err
	// }

	st.signature = txResp
	publisher.PublishMessage(st.Task, fmt.Sprintf("Tx Sent: %s", txResp))
	return nil
}

func (st *Transaction) ConfirmTransaction(ctx context.Context, publisher subscriptionhub.Publisher) error {
	stream := make(chan client.ConfirmMessage, 100)

	go func(stream chan client.ConfirmMessage) {
		defer close(stream)
		client.ConfirmTransactionWithStream(ctx, st.getHttpClient(), st.signature, stream)
	}(stream)

	for msg := range stream {
		if msg.Err != "" {
			return fmt.Errorf("%v", msg.Err)
		}

		publisher.PublishMessage(st.Task, msg.Message)
	}

	return nil
}

func (st *Transaction) GetSignature() string {
	return st.signature.String()
}

func (st *Transaction) UpdatePosition(ctx context.Context, publisher subscriptionhub.Publisher) (tokenAmount, solAmount float64, pos *positionmodel.Position, err error) {
	tokenAmount, solAmount, err = st.extractTokenAndSolFromTx(ctx, st.signature)
	if err != nil {
		return
	}

	tokensSold := new(big.Float).SetFloat64(tokenAmount)
	solReceived := new(big.Float).SetFloat64(solAmount)
	positionReturn := &positionmodel.Position{}

	if st.Task.Position_id == nil {
		position, exists := st.PositionService.FindPositionIfExists(st.Task.Token, st.Task.Wallet.PublicKey())
		if !exists {
			return tokenAmount, solAmount, pos, fmt.Errorf("position not found for sell task %d: ", st.Task.Id())
		}

		err = st.PositionService.ReportSell(position.PositionId, tokensSold, solReceived)
		if err != nil {
			return
		}

		positionReturn = position
	} else {
		err = st.PositionService.ReportSell(*st.Task.Position_id, tokensSold, solReceived)
		if err != nil {
			return
		}

		pos, err = st.PositionService.GetById(*st.Task.Position_id)
		if err != nil {
			return tokenAmount, solAmount, positionReturn, fmt.Errorf("position not found for sell task %d: ", st.Task.Id())
		}
		positionReturn = pos

	}
	return tokenAmount, solAmount, positionReturn, nil
}

func (st *Transaction) extractTokenAndSolFromTx(ctx context.Context, signature solana.Signature) (tokenAmount float64, solAmount float64, err error) {
	tx, err := st.Task.HttpClient().GetParsedTransaction(ctx, signature, &rpc.GetParsedTransactionOpts{Commitment: rpc.CommitmentConfirmed, MaxSupportedTransactionVersion: &rpc.MaxSupportedTransactionVersion0})
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
	walletIndex := -1

	for i, account := range transactionMessage.AccountKeys {
		if account.PublicKey == walletPubkey {
			walletIndex = i
		}
	}

	if walletIndex == -1 {
		return tokenAmount, solAmount, fmt.Errorf("could not find user's wallet in account keys")
	}

	preBalance := int64(tx.Meta.PreBalances[walletIndex])
	postBalance := int64(tx.Meta.PostBalances[walletIndex])

	solAmountLamport := postBalance - preBalance
	solAmount = float64(solAmountLamport)

	return tokenAmount, solAmount, nil
}

func (st *Transaction) GetTask() tasks.Task {
	return st.Task
}

func (st *Transaction) getHttpClient() *rpc.Client {
	return st.Task.HttpClient()
}

func (st *Transaction) getAllInstructionsForSell(ctx context.Context, sellTask *tasks.SellTask) ([]solana.Instruction, error) {
	// we need to check if the user passed a positionId -> if so position == nil
	// else we retrieve the positonId
	var pos *positionmodel.Position
	if sellTask.Position_id != nil {
		position, err := st.PositionService.GetById(*sellTask.Position_id)
		if err != nil {
			return nil, err
		}

		isTokensRemaining := position.TokenRemaining.Cmp(big.NewFloat(0)) == 1
		if !isTokensRemaining {
			return nil, fmt.Errorf("no tokens remaining")
		}

		pos = position
	} else {
		//report fake buy for positions
		ata, err := client.GetATA(ctx, sellTask.Wallet.PublicKey(), sellTask.Token, sellTask.HttpClient())
		if err != nil {
			return nil, err
		}

		tokenAmount, err := client.GetTokenAccountBalance(ctx, ata, sellTask.HttpClient())
		if err != nil {
			return nil, err
		}

		tokenAmountBig := new(big.Float).SetUint64(*tokenAmount)
		err = st.PositionService.ReportBuy(ctx, sellTask.Id(), sellTask.Token, sellTask.Wallet.PublicKey(), tokenAmountBig, big.NewFloat(0))
		if err != nil {
			return nil, err
		}

		pos = nil
	}

	computeLimitInstruction := instructions.GetComputeUnitLimitInstruction(sellTask.ComputeUnits)
	computeLimitBudgetInstruction := instructions.GetComputeUnitBudgetInstruction(sellTask.Fee, sellTask.ComputeUnits)

	sellInstructions, err := pumpInstructions.GetSellInstruction(ctx, sellTask, pos)
	if err != nil {
		logger.Error("Error getting sell instruction", err)
		return nil, err
	}

	solInstructions := []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, sellInstructions}
	return solInstructions, nil
}
