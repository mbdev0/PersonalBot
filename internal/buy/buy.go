// connect to http server
// create transaction
// send transaction to http server
// or can we do this via jito?
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

	txResp, err := client.SimulateTransaction(context.Background(), tx)
	if err != nil {
		logger.Error(err)
	}

	fmt.Println(txResp.Value)
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

	instructions := []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, idEmponenetInstruction, buyInstruction}

	return instructions

}
