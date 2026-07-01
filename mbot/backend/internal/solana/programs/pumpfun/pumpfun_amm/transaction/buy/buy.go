package buy

import (
	"context"
	"fmt"
	"math/big"
	"personal_bot/backend/infrastructure/solana_price"
	"personal_bot/backend/internal/core/constants"
	positionmodels "personal_bot/backend/internal/core/position"
	"personal_bot/backend/internal/core/tasks"
	subscriptionhub "personal_bot/backend/internal/services/subscription_hub"
	"personal_bot/backend/internal/solana/instructions"
	amminstructions "personal_bot/backend/internal/solana/programs/pumpfun/pumpfun_amm/instructions"
	"personal_bot/backend/internal/solana/programs/pumpfun/pumpfun_amm/models"
	"personal_bot/backend/internal/solana/programs/pumpfun/pumpfun_amm/pda"
	ammtransactiondecoder "personal_bot/backend/internal/solana/programs/pumpfun/pumpfun_amm/transaction_decoder"
	"personal_bot/backend/internal/solana/transaction"

	transactiondecoder "personal_bot/backend/internal/solana/programs/pumpfun/pumpfun_native/transaction_decoder"

	"personal_bot/backend/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type Transaction struct {
	BuyTask         *tasks.BuyTask
	instructions    *[]solana.Instruction
	transaction     *solana.Transaction
	signature       solana.Signature
	PositionService positionmodels.PositionService
	poolAddress     string
	Publisher       subscriptionhub.Publisher
}

func NewTransaction(task *tasks.BuyTask, posService positionmodels.PositionService, publisher subscriptionhub.Publisher) *Transaction {
	return &Transaction{
		BuyTask:         task,
		PositionService: posService,
		Publisher:       publisher,
	}
}

func (bt *Transaction) BuildInstructions(ctx context.Context) error {
	if bt.BuyTask == nil {
		return fmt.Errorf("buy task was nil - make sure buy task is set")
	}

	buyInstructions, err := bt.getAllInstructionsForBuy(ctx, bt.BuyTask)
	if buyInstructions == nil || err != nil {
		logger.Error("Error creating buy instructions - no instructions created")
		return err
	}

	bt.instructions = &buyInstructions
	bt.Publisher.PublishMessage(bt.BuyTask, "Instructions Built")

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

func (bt *Transaction) BuildTransaction(ctx context.Context) error {
	if bt.BuyTask == nil {
		return fmt.Errorf("buy task was nil - make sure buy task is set")
	}

	tx, err := transaction.BuildTx(ctx, bt.getHttpClient(), bt.BuyTask.Wallet, bt.instructions)
	if err != nil {
		return err
	}

	bt.transaction = tx
	bt.Publisher.PublishMessage(bt.BuyTask, "TX Built")
	return nil
}

func (bt *Transaction) SendTransaction(ctx context.Context) error {
	rpcClient := bt.getHttpClient()
	signature, err := transaction.SendTx(ctx, rpcClient, bt.transaction)
	if err != nil {
		return err
	}

	bt.signature = signature
	bt.Publisher.PublishMessage(bt.BuyTask, fmt.Sprintf("Tx Sent: %s", signature))
	return nil
}

func (bt *Transaction) ConfirmTransaction(ctx context.Context) error {
	rpcClient := bt.getHttpClient()
	return transaction.ConfirmTx(ctx, rpcClient, bt.signature, bt.BuyTask, bt.Publisher)
}

func (bt *Transaction) UpdatePosition(ctx context.Context) (tokenAmount, solAmount float64, pos *positionmodels.Position, err error) {
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

	pool, err := pda.GetPoolPda(bt.BuyTask.Token.String(), constants.WSOLTokenAddress)
	if err != nil {
		return
	}

	payload := positionmodels.ReportBuyPayload{
		BuyTaskId:            bt.BuyTask.Id(),
		StrategyId:           bt.BuyTask.StrategyId,
		TokenAddress:         bt.BuyTask.Token,
		WalletAddress:        bt.BuyTask.Wallet.PublicKey(),
		TokenAmount:          new(big.Float).SetFloat64(float64(buyEvent.BaseAmountOut)),
		SolSpent:             new(big.Float).SetFloat64(solAmnt),
		MarketCap:            marketCapUSD,
		AddressForUrl:        bt.poolAddress,
		AddressForMonitoring: pool,
		Program:              bt.BuyTask.Program(),
		Logger:               bt.BuyTask.Logger(),
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

func (bt *Transaction) GetSignature() string {
	return bt.signature.String()
}

func (bt *Transaction) GetAddressForURL() (string, error) {
	return bt.poolAddress, nil
}
