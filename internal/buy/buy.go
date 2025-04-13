package buy

import (
	"context"
	"fmt"
	"math/big"
	"pump_fun/internal/handlers"
	instructionbuilder "pump_fun/internal/instruction_builder"
	"pump_fun/internal/logger"
	rpcclient "pump_fun/internal/rpc_client"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func SendBuyTransaction(privateKey solana.PrivateKey, tokenAddress string, buyAmountLamport big.Int, slippage float64) {
	sendBuyTransaction(privateKey, tokenAddress, buyAmountLamport, slippage)
}

func sendBuyTransaction(privateKey solana.PrivateKey, tokenAddress string, buyAmountLamport big.Int, slippage float64) {
	instructions := getAllInstructionsForBuy(privateKey, tokenAddress, buyAmountLamport, slippage)
	client := rpcclient.GetClient()

	latestHash, err := client.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		logger.Error(err)
	}

	tx, err := solana.NewTransaction(instructions, latestHash.Value.Blockhash, solana.TransactionPayer(privateKey.PublicKey()))
	if err != nil {
		logger.Error(err)
	}

	handlers.SignTx(tx, privateKey)

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

func getAllInstructionsForBuy(privateKey solana.PrivateKey, tokenAddress string, buyAmountLamport big.Int, slippage float64) (buyInstructions []solana.Instruction) {
	tokenAddressPubkey, err := solana.PublicKeyFromBase58(tokenAddress)
	if err != nil {
		logger.Error(err)
	}

	computeLimitInstruction := instructionbuilder.GetComputeUnitLimitInstruction(80000)
	computeLimitBudgetInstruction := instructionbuilder.GetComputeUnitBudgetInstruction(500000)
	idEmponenetInstruction := instructionbuilder.GetIdEmponentInstruction(privateKey.PublicKey(), tokenAddressPubkey)
	buyInstruction, err := instructionbuilder.GetBuyInstruction(tokenAddress, privateKey.PublicKey().String(), buyAmountLamport, slippage)

	if err != nil {
		logger.Error(err)
	}

	if idEmponenetInstruction != nil {
		return []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, idEmponenetInstruction, buyInstruction}
	} else {
		return []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, buyInstruction}
	}
}
