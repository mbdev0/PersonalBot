package buy

import (
	"context"
	"fmt"
	"math/big"
	lookuptable "personal_bot/app/lookup_table"
	"personal_bot/infrastructure/solana_price"
	"personal_bot/internal/core/constants"
	positionmodel "personal_bot/internal/core/position"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/services/position"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	"personal_bot/internal/solana/client"
	"personal_bot/internal/solana/instructions"
	ammPda "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/pda"

	amminstructions "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/instructions"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/models"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/pool"
	ammtransactiondecoder "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/transaction_decoder"
	transactiondecoder "personal_bot/internal/solana/programs/pumpfun/pumpfun_native/transaction_decoder"
	wallets "personal_bot/internal/solana/wallet"
	"personal_bot/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type Transaction struct {
	SellTask        *tasks.SellTask
	instructions    *[]solana.Instruction
	transaction     *solana.Transaction
	signature       solana.Signature
	PositionService *position.Service
	poolAddress     string
}

func (st *Transaction) BuildInstructions(ctx context.Context, publisher subscriptionhub.Publisher) error {
	if st.SellTask == nil {
		return fmt.Errorf("sell task was nil - make sure buy task is set")
	}

	sellInstructions, err := st.getAllInstructionsForSell(ctx, st.SellTask)
	if sellInstructions == nil || err != nil {
		logger.Error("Error creating buy instructions - no instructions created")
		return err
	}

	st.instructions = &sellInstructions
	publisher.PublishMessage(st.SellTask, "Instructions Built")

	return nil
}

func (st *Transaction) getAllInstructionsForSell(ctx context.Context, sellTask *tasks.SellTask) (sellInstructions []solana.Instruction, err error) {
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

		//pool address, token address, tokenAccount -> if its new or not

		//thing is if we want to check for the correct token its ANOTHER account call - but for standalone sells does this matter?
		poolAddy, err := ammPda.GetPoolPda(sellTask.Token.String(), constants.WSOLTokenAddress)
		if err != nil {
			return nil, err
		}

		isNewToken, err := instructions.IsTokenAccountNew(ctx, solana.MustPublicKeyFromBase58(poolAddy), st.getHttpClient())
		if err != nil {
			return nil, err
		}

		baseMintAddy, err := ammPda.GetPoolTokenAccount(poolAddy, sellTask.Token.String(), isNewToken)
		if err != nil {
			return nil, err
		}

		quoteMintAddy, err := ammPda.GetPoolTokenAccount(poolAddy, constants.WSOLTokenAddress, isNewToken)
		if err != nil {
			return nil, err
		}

		poolBalances, err := pool.GetTokenBalances(ctx, baseMintAddy, quoteMintAddy, st.getHttpClient())
		if err != nil {
			return nil, err
		}

		marketCap, err := pool.GetMarketCapUSD(ctx, poolBalances)

		// marketCap, err, _ := bondingcurve.GetMarketCapFromTokenAddress(ctx, st.SellTask.Token, st.SellTask.HttpClient())
		if err != nil {
			return nil, err
		}

		//TODO: add address for URL
		payload := positionmodel.ReportBuyPayload{
			BuyTaskId:     st.SellTask.Id(),
			StrategyId:    st.SellTask.StrategyId,
			TokenAddress:  st.SellTask.Token,
			WalletAddress: st.SellTask.Wallet.PublicKey(),
			TokenAmount:   tokenAmountBig,
			SolSpent:      new(big.Float).SetFloat64(0),
			MarketCap:     marketCap,
			AddressForUrl: poolAddy,
		}

		err = st.PositionService.ReportBuy(ctx, payload)
		if err != nil {
			return nil, err
		}

		pos = nil
	}

	computeLimitInstruction := instructions.GetComputeUnitLimitInstruction(sellTask.ComputeUnits)
	computeLimitBudgetInstruction := instructions.GetComputeUnitBudgetInstruction(sellTask.Fee, sellTask.ComputeUnits)
	// idEmponenetInstruction, err := instructions.GetIdempotentInstruction(ctx, sellTask.Wallet.PublicKey(), sellTask.Token, sellTask.HttpClient())
	// if err != nil {
	// 	return
	// }

	// wsolAddress := solana.MustPublicKeyFromBase58(constants.WSOLTokenAddress)
	// idEmponenetInstructionWsol, err := instructions.GetIdempotentInstruction(ctx, sellTask.Wallet.PublicKey(), wsolAddress, sellTask.HttpClient())
	// if err != nil {
	// 	return
	// }

	sellInstruction, err := amminstructions.GetSellInstruction(ctx, sellTask, pos)
	if err != nil {
		logger.Error("Error creating buy instruction", err)
		return
	}

	poolAccount := sellInstruction.AccountValues.Get(0)
	if poolAccount != nil {
		st.poolAddress = poolAccount.PublicKey.String()
	}

	// sellInstructions := []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, sellInstruction}
	// if idEmponenetInstruction != nil {
	// 	buyInstructions = append(buyInstructions, idEmponenetInstruction)
	// }

	// if idEmponenetInstructionWsol != nil {
	// 	buyInstructions = append(buyInstructions, idEmponenetInstructionWsol)
	// }

	// buyInstructions = append(buyInstructions, sellInstruction)

	return []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, sellInstruction}, nil
}

func (st *Transaction) BuildTransaction(ctx context.Context, publisher subscriptionhub.Publisher) error {
	if st.SellTask == nil {
		return fmt.Errorf("sell task was nil - make sure sell task is set")
	}

	latestHash, err := client.GetLatestBlockhash(ctx, st.getHttpClient())
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
		return err
	}

	opts := []solana.TransactionOption{
		solana.TransactionPayer(st.SellTask.Wallet.PublicKey()),
	}

	//TODO - Add constants into our address lookup table
	accountLookupMap, err := lookuptable.GetAddressLookupTable(st.getHttpClient())
	if err != nil {
		logger.Error("Error getting address lookup table, proceeding without it: ", err)
	} else {
		opts = append(opts, solana.TransactionAddressTables(accountLookupMap))
	}

	tx, err := solana.NewTransaction(*st.instructions, latestHash.Value.Blockhash, opts...)
	if err != nil {
		logger.Error("Error creating transaction", err)
		return err
	}

	tx.Message.SetVersion(solana.MessageVersionV0)

	err = wallets.SignTx(tx, st.SellTask.Wallet)
	if err != nil {
		return err
	}
	st.transaction = tx
	publisher.PublishMessage(st.SellTask, "TX Built")

	return nil
}

func (st *Transaction) SendTransaction(ctx context.Context, publisher subscriptionhub.Publisher) error {
	rpcClient := st.getHttpClient()
	// SIMULATE TRANSACTION
	// txResp, err := rpcClient.SimulateTransaction(bt.BuyTask.Ctx(), bt.transaction)
	// if err != nil {
	// 	logger.Error("Transaction simulation failed", err)
	// 	return nil
	// }
	// fmt.Println(txResp.Value)

	// SEND TRANSACTION WITH OPTIONS
	// maxRetries := uint(5)
	txResp, err := rpcClient.SendTransactionWithOpts(ctx, st.transaction, rpc.TransactionOpts{Encoding: solana.EncodingBase64, SkipPreflight: true})
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

	st.signature = txResp
	publisher.PublishMessage(st.SellTask, fmt.Sprintf("Tx Sent: %s", txResp))
	return nil
}

func (st *Transaction) ConfirmTransaction(ctx context.Context, publisher subscriptionhub.Publisher) error {
	rpcClient := st.getHttpClient()

	stream := make(chan client.ConfirmMessage, 100)

	go func(stream chan client.ConfirmMessage) {
		defer close(stream)
		client.ConfirmTransactionWithStream(ctx, rpcClient, st.signature, stream)
	}(stream)

	for msg := range stream {
		if msg.Err != "" {
			publisher.PublishMessage(st.SellTask, msg.Err)
			return fmt.Errorf("%v", msg.Err)
		}
		publisher.PublishMessage(st.SellTask, msg.Message)
	}

	return nil
}

func (st *Transaction) UpdatePosition(ctx context.Context, publisher subscriptionhub.Publisher) (tokenAmount, solAmount float64, pos *positionmodel.Position, err error) {
	solAmnt, sellEvent, err := st.extractTokenAndSolFromTx(ctx, st.signature)
	if err != nil {
		return
	}

	solPrice, err := solana_price.GetSolPrice()
	if err != nil {
		return
	}

	pricePerToken := (float64(sellEvent.PoolQuoteTokenReserves) / constants.LamportsConversion) /
		(float64(sellEvent.PoolBaseTokenReserves) / constants.TokenAmountDecimals)

	marketCapUSD := new(big.Float).SetFloat64((pricePerToken * 1_000_000_000) * *solPrice)

	if st.SellTask.Position_id != nil {
		err = st.PositionService.ReportSell(ctx, *st.SellTask.Position_id, big.NewFloat(float64(sellEvent.UserQuoteAmountOut)), big.NewFloat(solAmnt), marketCapUSD)
		if err != nil {
			return
		}
	} else {
		err = st.PositionService.ReportSell(ctx, st.SellTask.Id(), big.NewFloat(float64(sellEvent.UserQuoteAmountOut)), big.NewFloat(solAmnt), marketCapUSD)
		if err != nil {
			return
		}

	}

	return float64(sellEvent.UserQuoteAmountOut), solAmnt, nil, nil
}

func (st *Transaction) extractTokenAndSolFromTx(ctx context.Context, signature solana.Signature) (solAmount float64, sellEvent models.SellEvent, err error) {
	solClient := st.getHttpClient()
	tx, err := solClient.GetParsedTransaction(ctx, signature, &rpc.GetParsedTransactionOpts{Commitment: rpc.CommitmentConfirmed, MaxSupportedTransactionVersion: &rpc.MaxSupportedTransactionVersion0})
	if err != nil {
		return
	}

	if tx.Meta.Err != nil {
		return solAmount, sellEvent, fmt.Errorf("error in transaction whilst extracting token amount + sol amount")
	}

	sellEvent, err = ammtransactiondecoder.GetSellEvent(tx)
	if err != nil {
		return
	}

	solAmount, err = transactiondecoder.ExtractTotalSolSpent(tx, st.SellTask.Wallet.PublicKey())
	if err != nil {
		return
	}

	return -solAmount, sellEvent, nil
}

func (st *Transaction) getHttpClient() *rpc.Client {
	return st.SellTask.HttpClient()
}

func (st *Transaction) GetTask() tasks.Task {
	return st.SellTask
}

func (st *Transaction) GetSignature() string {
	return st.signature.String()
}

func (st *Transaction) GetAddressForURL() (string, error) {
	return st.poolAddress, nil
}
