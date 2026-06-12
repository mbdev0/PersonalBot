package sell

import (
	"context"
	"fmt"
	"math/big"
	"personal_bot/infrastructure/solana_price"
	"personal_bot/internal/core/constants"
	positionmodel "personal_bot/internal/core/position"
	"personal_bot/internal/core/tasks"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	"personal_bot/internal/solana/client"
	"personal_bot/internal/solana/instructions"
	bondingcurve "personal_bot/internal/solana/programs/pumpfun/bonding_curve"
	pumpInstructions "personal_bot/internal/solana/programs/pumpfun/pumpfun_native/instructions"
	pumpmodels "personal_bot/internal/solana/programs/pumpfun/pumpfun_native/models"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_native/pda"
	transactiondecoder "personal_bot/internal/solana/programs/pumpfun/pumpfun_native/transaction_decoder"
	"personal_bot/internal/solana/transaction"

	"personal_bot/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type Transaction struct {
	Task            *tasks.SellTask
	instructions    []solana.Instruction
	transaction     *solana.Transaction
	signature       solana.Signature
	PositionService positionmodel.PositionService
	Publisher       subscriptionhub.Publisher
}

func NewTransaction(task *tasks.SellTask, posService positionmodel.PositionService, publisher subscriptionhub.Publisher) *Transaction {
	return &Transaction{
		Task:            task,
		PositionService: posService,
		Publisher:       publisher,
	}
}

func (st *Transaction) BuildInstructions(ctx context.Context) error {

	sellInstructions, err := st.getAllInstructionsForSell(ctx, st.Task)
	if err != nil {
		logger.Error("Error getting instructions for sell task", err)
		return err
	}

	if len(sellInstructions) == 0 {
		return fmt.Errorf("sell instruction's weren't generated properly - check sell instruction builder")
	}

	st.instructions = sellInstructions
	st.Publisher.PublishMessage(st.Task, "instructions built")
	return nil
}

func (st *Transaction) BuildTransaction(ctx context.Context) error {
	if st.Task == nil {
		return fmt.Errorf("sell task is null - check if sell task was set")
	}

	tx, err := transaction.BuildTx(ctx, st.getHttpClient(), st.Task.Wallet, &st.instructions)
	if err != nil {
		return err
	}

	st.transaction = tx
	st.Publisher.PublishMessage(st.Task, "tx built")

	return nil
}

func (st *Transaction) SendTransaction(ctx context.Context) error {
	rpcClient := st.getHttpClient()
	sig, err := transaction.SendTx(ctx, rpcClient, st.transaction)
	if err != nil {
		return err
	}

	st.signature = sig
	st.Publisher.PublishMessage(st.Task, fmt.Sprintf("Tx Sent: %s", sig))
	return nil
}

func (st *Transaction) ConfirmTransaction(ctx context.Context) error {
	return transaction.ConfirmTx(ctx, st.getHttpClient(), st.signature, st.GetTask(), st.Publisher)
}

func (st *Transaction) GetSignature() string {
	return st.signature.String()
}

func (st *Transaction) UpdatePosition(ctx context.Context) (tokenAmount, solAmount float64, pos *positionmodel.Position, err error) {
	solAmount, tradeEvent, err := st.extractTokenAndSolFromTx(ctx, st.signature)
	if err != nil {
		return
	}

	tokenAmount = float64(tradeEvent.TokenAmount)
	logger.Information("token amount: ", tokenAmount)

	solPrice, err := solana_price.GetSolPrice()
	if err != nil {
		return
	}

	pricePerToken := (float64(tradeEvent.VirtualSolReserves) / constants.LamportsConversion) /
		(float64(tradeEvent.VirtualTokenReserves) / constants.TokenAmountDecimals)

	marketCapUSD := new(big.Float).SetFloat64((pricePerToken * 1_000_000_000) * *solPrice)

	tokensSold := new(big.Float).SetFloat64(tokenAmount)
	solReceived := new(big.Float).SetFloat64(solAmount)

	if st.Task.Position_id == nil {
		position, exists := st.PositionService.FindPositionIfExists(st.Task.Token, st.Task.Wallet.PublicKey())
		if !exists {
			return 0, 0, nil, fmt.Errorf("position not found for sell task %d", st.Task.Id())
		}
		if err = st.PositionService.ReportSell(ctx, position.PositionId, tokensSold, solReceived, marketCapUSD); err != nil {
			return
		}
		return tokenAmount, solAmount, position, nil
	}

	if err = st.PositionService.ReportSell(ctx, *st.Task.Position_id, tokensSold, solReceived, marketCapUSD); err != nil {
		return
	}

	pos, err = st.PositionService.GetById(*st.Task.Position_id)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("position not found for sell task %d", st.Task.Id())
	}

	return tokenAmount, solAmount, pos, nil
}

func (st *Transaction) extractTokenAndSolFromTx(ctx context.Context, signature solana.Signature) (solAmount float64, tradeEvent pumpmodels.TradeEvent, err error) {
	tx, err := st.Task.HttpClient().GetParsedTransaction(ctx, signature, &rpc.GetParsedTransactionOpts{Commitment: rpc.CommitmentConfirmed, MaxSupportedTransactionVersion: &rpc.MaxSupportedTransactionVersion0})
	if err != nil {
		return
	}

	if tx.Meta.Err != nil {
		return solAmount, tradeEvent, fmt.Errorf("error in transaction whilst extracting token amount + sol amount")
	}

	solAmount, err = transactiondecoder.ExtractTotalSolSpent(tx, st.Task.Wallet.PublicKey())
	if err != nil {
		return
	}

	tradeEvent, err = transactiondecoder.GetTradeEvent(tx)
	if err != nil {
		return
	}

	return
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
		marketCap, err, _ := bondingcurve.GetMarketCapFromTokenAddress(ctx, st.Task.Token, st.Task.HttpClient())
		if err != nil {
			return nil, err
		}

		bondingCurve, err := pda.GetBondingCurveAddress(st.Task.GetToken())
		if err != nil {
			return nil, err
		}

		payload := positionmodel.ReportBuyPayload{
			BuyTaskId:     st.Task.Id(),
			StrategyId:    st.Task.StrategyId,
			TokenAddress:  st.Task.Token,
			WalletAddress: st.Task.Wallet.PublicKey(),
			TokenAmount:   tokenAmountBig,
			SolSpent:      new(big.Float).SetFloat64(0),
			MarketCap:     marketCap,
			AddressForUrl: bondingCurve,
			Program:       st.Task.Program(),
		}

		err = st.PositionService.ReportBuy(ctx, payload)
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

func (st *Transaction) GetAddressForURL() (string, error) {
	return pda.GetBondingCurveAddress(st.Task.GetToken())
}
