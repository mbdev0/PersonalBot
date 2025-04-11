// connect to http server
// create transaction
// send transaction to http server
// or can we do this via jito?
package buy

import (
	"context"
	"fmt"
	"math/big"
	instructionbuilder "pump_fun/internal/instruction_builder"
	rpcclient "pump_fun/internal/rpc_client"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func SendBuyTransaction(account string, tokenAddress string, buyAmountLamport big.Int, slippage float64) {
	sendBuyTransaction(account, tokenAddress, buyAmountLamport, slippage)
}

func sendBuyTransaction(account string, tokenAddress string, buyAmountLamport big.Int, slippage float64) {

	privateKey, err := solana.PrivateKeyFromBase58(account)

	if err != nil {
		fmt.Println(err)
	}

	tokenAddressPubkey, err := solana.PublicKeyFromBase58(tokenAddress)
	if err != nil {
		fmt.Println(err)
	}

	computeLimitInstruction := instructionbuilder.GetComputeUnitLimitInstruction(80000)
	computeLimitBudgetInstruction := instructionbuilder.GetComputeUnitBudgetInstruction(500000)
	idEmponenetInstruction := instructionbuilder.GetIdEmponentInstruction(privateKey.PublicKey(), tokenAddressPubkey)
	buyInstruction, err := instructionbuilder.GetBuyInstruction(tokenAddress, privateKey.PublicKey().String(), buyAmountLamport, slippage)

	if err != nil {
		return
	}
	client := rpcclient.GetClient()

	instructions := []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, idEmponenetInstruction, buyInstruction}
	latestHash, err := client.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)

	if err != nil {
		fmt.Println("ERROR", err)
	}

	tx, err := solana.NewTransaction(instructions, latestHash.Value.Blockhash, solana.TransactionPayer(privateKey.PublicKey()))

	if err != nil {
		fmt.Println("ERROR", err)
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if privateKey.PublicKey().Equals(key) {
				return &privateKey
			}
			return nil
		},
	)

	if err != nil {
		panic(fmt.Errorf("unable to sign transaction: %w", err))
	}

	txResp, err := client.SimulateTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(txResp.Value)

}
