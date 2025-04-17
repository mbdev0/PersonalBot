package buy

import (
	"context"
	"fmt"
	"math/big"
	"pump_fun/internal/handlers"
	instructionbuilder "pump_fun/internal/instruction_builder"
	"pump_fun/internal/logger"
	"pump_fun/internal/models/tasks"
	rpcclient "pump_fun/internal/rpc_client"

	"github.com/gagliardetto/solana-go"
)

func SendBuyTransaction(buyTask *tasks.BuyTask) {
	sendBuyTransaction(buyTask)
}

func sendBuyTransaction(buyTask *tasks.BuyTask) {
	instructions := getAllInstructionsForBuy(buyTask.Wallet, buyTask.TokenAddress, buyTask.BuyAmount, buyTask.Slippage)

	latestHash, err := rpcclient.GetLatestBlockhash()
	if err != nil {
		logger.Error(err)
	}

	tx, err := solana.NewTransaction(instructions, latestHash.Value.Blockhash, solana.TransactionPayer(buyTask.Wallet.PublicKey()))
	if err != nil {
		logger.Error(err)
	}

	handlers.SignTx(tx, buyTask.Wallet)

	client := rpcclient.GetClient()
	// SIMULATE TRANSACTION
	txResp, err := client.SimulateTransaction(context.Background(), tx)
	if err != nil {
		logger.Error(err)
	}
	fmt.Println(txResp.Value)

	// SEND TRANSACTION WITH OPTIONS
	// maxRetries := uint(5)
	// txResp, err := client.SendTransactionWithOpts(context.Background(), tx, rpc.TransactionOpts{Encoding: solana.EncodingBase64, SkipPreflight: false, MaxRetries: &maxRetries})
	// if err != nil {
	// 	logger.Error(err)
	// }
	// fmt.Println(txResp.String())

	// SEND TRANSACTION WITH NO OPTS
	// txResp, err := client.SendTransaction(context.Background(), tx)
	// if err != nil {
	// 	logger.Error(err)
	// }
	// fmt.Println(txResp.String())

}

func getAllInstructionsForBuy(privateKey solana.PrivateKey, tokenAddress solana.PublicKey, buyAmountLamport big.Int, slippage float64) (buyInstructions []solana.Instruction) {

	computeLimitInstruction := instructionbuilder.GetComputeUnitLimitInstruction(80000)
	computeLimitBudgetInstruction := instructionbuilder.GetComputeUnitBudgetInstruction(500000)
	idEmponenetInstruction := instructionbuilder.GetIdEmponentInstruction(privateKey.PublicKey(), tokenAddress)
	buyInstruction, err := instructionbuilder.GetBuyInstruction(tokenAddress.String(), privateKey.PublicKey().String(), buyAmountLamport, slippage)

	if err != nil {
		logger.Error(err)
	}

	if idEmponenetInstruction != nil {
		return []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, idEmponenetInstruction, buyInstruction}
	} else {
		return []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, buyInstruction}
	}
}
