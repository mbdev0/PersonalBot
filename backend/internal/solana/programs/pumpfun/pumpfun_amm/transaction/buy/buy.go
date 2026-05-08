package buy

import (
	"context"
	"fmt"
	"math/big"
	lookuptable "personal_bot/app/lookup_table"
	"personal_bot/infrastructure/solana_price"
	"personal_bot/internal/core/constants"
	positionmodels "personal_bot/internal/core/position"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/services/position"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	"personal_bot/internal/solana/client"
	"personal_bot/internal/solana/instructions"
	amminstructions "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/instructions"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/models"
	ammtransactiondecoder "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/transaction_decoder"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_native/pda"
	transactiondecoder "personal_bot/internal/solana/programs/pumpfun/pumpfun_native/transaction_decoder"

	wallets "personal_bot/internal/solana/wallet"

	"personal_bot/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type Transaction struct {
	BuyTask         *tasks.BuyTask
	instructions    *[]solana.Instruction
	transaction     *solana.Transaction
	signature       solana.Signature
	PositionService *position.Service
	poolAddress     string
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

func (bt *Transaction) getAllInstructionsForBuy(ctx context.Context, buyTask *tasks.BuyTask) (buyInstructions []solana.Instruction, err error) {
	computeLimitInstruction := instructions.GetComputeUnitLimitInstruction(buyTask.ComputeUnits)
	computeLimitBudgetInstruction := instructions.GetComputeUnitBudgetInstruction(buyTask.Fee, buyTask.ComputeUnits)
	idEmponenetInstruction, err := instructions.GetIdempotentInstruction(ctx, buyTask.Wallet.PublicKey(), buyTask.Token, buyTask.HttpClient())
	if err != nil {
		return
	}

	wsolAddress := solana.MustPublicKeyFromBase58(constants.WSOLTokenAddress)
	idEmponenetInstructionWsol, err := instructions.GetIdempotentInstruction(ctx, buyTask.Wallet.PublicKey(), wsolAddress, buyTask.HttpClient())
	if err != nil {
		return
	}

	buyInstruction, err := amminstructions.GetBuyInstruction(ctx, buyTask)
	if err != nil {
		logger.Error("Error creating buy instruction", err)
		return
	}

	poolAccount := buyInstruction.AccountValues.Get(0)
	if poolAccount != nil {
		bt.poolAddress = poolAccount.PublicKey.String()
	}

	buyInstructions = []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction}
	if idEmponenetInstruction != nil {
		buyInstructions = append(buyInstructions, idEmponenetInstruction)
	}

	if idEmponenetInstructionWsol != nil {
		
		buyInstructions = append(buyInstructions, idEmponenetInstructionWsol)
	}

	buyInstructions = append(buyInstructions, buyInstruction)

	return
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

	//TODO - Add constants into our address lookup table
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

func (bt *Transaction) UpdatePosition(ctx context.Context, publisher subscriptionhub.Publisher) (tokenAmount, solAmount float64, pos *positionmodels.Position, err error) {
	solAmnt, buyEvent, err := bt.extractTokenAndSolFromTx(ctx, bt.signature)
	if err != nil {
		return
	}

	solPrice, err := solana_price.GetSolPrice()
	if err != nil {
		return
	}

	pricePerToken := (float64(buyEvent.PoolQuoteTokenReserves) / constants.LamportsConversion) /
		(float64(buyEvent.PoolBaseTokenReserves) / constants.TokenAmountDecimals)

	marketCapUSD := new(big.Float).SetFloat64((pricePerToken * 1_000_000_000) * *solPrice)

	bondingCurve, err := pda.GetBondingCurveAddress(bt.BuyTask.GetToken())
	if err != nil {
		return
	}

	payload := positionmodels.ReportBuyPayload{
		BuyTaskId:           bt.BuyTask.Id(),
		StrategyId:          bt.BuyTask.StrategyId,
		TokenAddress:        bt.BuyTask.Token,
		WalletAddress:       bt.BuyTask.Wallet.PublicKey(),
		TokenAmount:         new(big.Float).SetFloat64(float64(buyEvent.BaseAmountOut)),
		SolSpent:            new(big.Float).SetFloat64(solAmnt),
		MarketCap:           marketCapUSD,
		AddressForUrl:       bt.poolAddress,
		AdressForMonitoring: bondingCurve,
	}

	err = bt.PositionService.ReportBuy(ctx, payload)
	if err != nil {
		return
	}

	return float64(buyEvent.BaseAmountOut), solAmnt, nil, nil
}

func (bt *Transaction) extractTokenAndSolFromTx(ctx context.Context, signature solana.Signature) (solAmount float64, buyEvent models.BuyEvent, err error) {
	solClient := bt.getHttpClient()
	tx, err := solClient.GetParsedTransaction(ctx, signature, &rpc.GetParsedTransactionOpts{Commitment: rpc.CommitmentConfirmed, MaxSupportedTransactionVersion: &rpc.MaxSupportedTransactionVersion0})
	if err != nil {
		return
	}

	if tx.Meta.Err != nil {
		return solAmount, buyEvent, fmt.Errorf("error in transaction whilst extracting token amount + sol amount")
	}

	buyEvent, err = ammtransactiondecoder.GetBuyEvent(tx)
	if err != nil {
		return
	}

	solAmount, err = transactiondecoder.ExtractTotalSolSpent(tx, bt.BuyTask.Wallet.PublicKey())
	if err != nil {
		return
	}

	return -solAmount, buyEvent, nil
}

func (bt *Transaction) getHttpClient() *rpc.Client {
	return bt.BuyTask.HttpClient()
}

func (bt *Transaction) GetTask() tasks.Task {
	return bt.BuyTask
}

func (bt *Transaction) GetSignature() string {
	return bt.signature.String()
}

func (bt *Transaction) GetAddressForURL() (string, error) {
	return bt.poolAddress, nil
}
