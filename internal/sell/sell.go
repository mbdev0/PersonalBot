package sell

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"pump_fun/internal/constants"
	"pump_fun/internal/handlers"
	instructionbuilder "pump_fun/internal/instruction_builder"
	"pump_fun/internal/logger"
	"pump_fun/internal/models/tasks"
	"pump_fun/internal/monitoring/transactions/bonding_curve_decoder"
	"pump_fun/internal/monitoring/transactions/program_derived_address"
	rpcclient "pump_fun/internal/rpc_client"

	"github.com/gagliardetto/solana-go"
)

func SendSellTransaction(sellTask *tasks.SellTask) {
	sendSellTransaction(sellTask)
}

func sendSellTransaction(sellTask *tasks.SellTask) {
	// SetComputeUnitLimit
	computeLimitInstruction := instructionbuilder.GetComputeUnitLimitInstruction(100_000)
	// SetComputeUnitPrice
	computeLimitBudgetInstruction := instructionbuilder.GetComputeUnitBudgetInstruction(0.01, 100_000)

	fmt.Println("Compute Limit Instruction:", computeLimitInstruction)
	fmt.Println("Compute Budget Instruction:", computeLimitBudgetInstruction)

	// ##Transfer
	// TODO: Find how to get the user
	sourceAccount := sellTask.PublicKey
	fmt.Println("User Public Key:", sourceAccount)
	// TODO: Find how to get the source, Fetch this using the coin address?
	destinationAccount := "pfnXi2FdpFUUn6VyoxUohNyWk2Nup3ruguTgK8jaZaF"
	fmt.Println("Destination Key:", destinationAccount)

	// ##Sell
	// get all accounts for sell
	tokenAddress := sellTask.TokenAddress.String()
	walletAddress := sellTask.Wallet.PublicKey().String()

	bondingCurveAddress, err := program_derived_address.GetBondingCurveAddress(tokenAddress)
	if err != nil {
		logger.Error("Error getting bonding curve address:", err)
		return
	}

	associatedBondingCurveAddress, err := program_derived_address.GetAssociatedBondingCurveAddress(bondingCurveAddress, tokenAddress)
	if err != nil {
		logger.Error("Error getting associated bonding curve address:", err)
		return
	}

	//TODO: Make a function to fetch creatorAddress, looks like this is the same for buy
	bondingCurveData, err, _ := bonding_curve_decoder.GetBondingCurveDataFromAddress(bondingCurveAddress)
	if err != nil {
		logger.Error("Error getting bonding curve data:", err)
		return
	}
	creatorAddress, err := program_derived_address.GetCreatorVaultAddress(bondingCurveData.DevWallet.String())
	if err != nil {
		logger.Error("Error getting creator address:", err)
		return
	}

	accounts := []*solana.AccountMeta{
		//TODO: Move building this instruction inside instruction builder
		instructionbuilder.GetAccountMeta(constants.GlobalAccount, false, false),
		instructionbuilder.GetAccountMeta(constants.FeeRecipient, true, false),
		instructionbuilder.GetAccountMeta(tokenAddress, false, false),
		instructionbuilder.GetAccountMeta(bondingCurveAddress, true, false),
		instructionbuilder.GetAccountMeta(associatedBondingCurveAddress, true, false),
		instructionbuilder.GetAccountMeta("3VriQemVpwKTC1e4EbbmT7ZaNr9C7TfGbFmpWQgJxKyF", true, false), // This is the token account of the user
		instructionbuilder.GetAccountMeta(walletAddress, true, true),
		instructionbuilder.GetAccountMeta(solana.SystemProgramID.String(), false, false),
		instructionbuilder.GetAccountMeta(creatorAddress, true, false),
		instructionbuilder.GetAccountMeta(constants.TokenProgram, false, false),
		instructionbuilder.GetAccountMeta(constants.EventAuthority, false, false),
		instructionbuilder.GetAccountMeta(constants.Program, false, false),
	}
	fmt.Println("Accounts:", accounts)

	//Sell testing Data - Works with this data
	// rawDataHex := "33e685a4017f83adc0faa83e0800000047a80e0000000000"
	// instructionData, err := hex.DecodeString(rawDataHex)
	// fmt.Println("Instruction Data:", instructionData)

	// TODO: create sell transaction data
	buf := new(bytes.Buffer)

	discriminator := constants.SellInstructionDiscriminator
	if _, err := buf.Write(discriminator[:]); err != nil {
		logger.Error("Error writing discriminator to buffer", err)
		return
	}

	amount := uint64(26292925434)
	if err := binary.Write(buf, binary.LittleEndian, amount); err != nil {
		logger.Error("Error writing token amount", err)
		return
	}

	min_sol_output := uint64(582231)
	if err := binary.Write(buf, binary.LittleEndian, min_sol_output); err != nil {
		logger.Error("Error writing min_sol_output", err)
		return
	}

	fmt.Print("buf", buf.Bytes())

	// // NOT NEEDED BUT FEELS LIKE A WASTE TO THROW AWAY - Reading raw data
	// buf := bytes.NewBuffer(rawBytes)
	// type SellData struct {
	// 	Discriminator [8]byte
	// 	Amount        uint64
	// 	MinSolOutput  uint64
	// }
	// var sell SellData
	// binary.Read(buf, binary.LittleEndian, &sell)
	// fmt.Printf("Sell Amount: %d, Min SOL Output: %d\n", sell.Amount, sell.MinSolOutput)
	// fmt.Printf("Discriminator: %d\n", sell.Discriminator)

	// get latest block hash
	latestHash, err := rpcclient.GetLatestBlockhash()
	if err != nil {
		logger.Error("Error getting latest blockhash", err)
		return
	}

	// create the transaction
	sell_instructions := solana.NewInstruction(solana.MustPublicKeyFromBase58(constants.Program), accounts, instructionData)
	instructions := []solana.Instruction{computeLimitInstruction, computeLimitBudgetInstruction, sell_instructions}
	tx, err := solana.NewTransaction(instructions, latestHash.Value.Blockhash, solana.TransactionPayer(sellTask.Wallet.PublicKey()))
	if err != nil {
		logger.Error("Error creating transaction", err)
		return
	}

	// sign the transaction
	handlers.SignTx(tx, sellTask.Wallet)

	// simulate the transaction
	client := rpcclient.GetClient()
	txResp, err := client.SimulateTransaction(context.Background(), tx)
	if err != nil {
		logger.Error("Transaction simulation failed", err)
		return
	}
	fmt.Println(txResp.Value)
}
