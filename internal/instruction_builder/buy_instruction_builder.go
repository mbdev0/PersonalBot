package instructionbuilder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"pump_fun/internal/constants"
	"pump_fun/internal/handlers"
	"pump_fun/internal/monitoring/transactions/bonding_curve_decoder"

	"github.com/gagliardetto/solana-go"
)

func GetBuyInstruction(tokenAddress string, walletAddress string, sol_lamport_buy_amount big.Int, slippage float64) (instruction *solana.GenericInstruction, err error) {

	bondingCurveAddress, err := bonding_curve_decoder.GetBondingCurveAddress(tokenAddress)

	if err != nil {
		return nil, err
	}

	associatedBondingCurveAddress, err := bonding_curve_decoder.GetAssociatedBondingCurveAddress(bondingCurveAddress, tokenAddress)
	fmt.Println(associatedBondingCurveAddress)

	if err != nil {
		return nil, err
	}

	associatedTokenAddressPubkey, _, err := solana.FindAssociatedTokenAddress(solana.MustPublicKeyFromBase58(walletAddress), solana.MustPublicKeyFromBase58(tokenAddress))

	fmt.Println(walletAddress)
	fmt.Println(associatedTokenAddressPubkey)

	if err != nil {
		return nil, err
	}

	accounts := []*solana.AccountMeta{
		GetAccountMeta(constants.GlobalAccount, true, false),
		GetAccountMeta(constants.FeeRecipient, true, false),
		GetAccountMeta(tokenAddress, false, false),
		GetAccountMeta(bondingCurveAddress, true, false),
		GetAccountMeta(associatedBondingCurveAddress, true, false),
		GetAccountMeta(associatedTokenAddressPubkey.String(), true, false),
		GetAccountMeta(walletAddress, true, true),
		GetAccountMeta(solana.SystemProgramID.String(), false, false),
		GetAccountMeta(constants.TokenProgram, false, false),
		GetAccountMeta("SysvarRent111111111111111111111111111111111", false, false),
		GetAccountMeta(constants.EventAuthority, false, false),
		GetAccountMeta(constants.Program, false, false),
	}

	instruction_data := create_buy_data(sol_lamport_buy_amount, bondingCurveAddress, slippage)

	buy_instructions := solana.NewInstruction(solana.MustPublicKeyFromBase58(constants.Program), accounts, instruction_data)

	return buy_instructions, nil

}

func create_buy_data(sol_lamport_buy_amount big.Int, bondingCurveAddr string, slippage float64) []byte {
	// Get the token amount
	tokenAmount, err, hasCompleted := handlers.GetBuyTokenAmountFrom(sol_lamport_buy_amount, bondingCurveAddr)
	if err != nil || hasCompleted {
		fmt.Println(err, hasCompleted)
		return nil
	}

	// Create a buffer with capacity for:
	// 8 bytes discriminator + 8 bytes tokenAmount + 8 bytes sol_lamport_buy_amount
	buf := new(bytes.Buffer)

	// 1. Write the 8-byte discriminator (exact format depends on its type)
	discriminator := constants.BuyInstructionDiscriminator
	if _, err := buf.Write(discriminator[:]); err != nil {
		fmt.Println("Error writing discriminator:", err)
		return nil
	}

	// 2. Write tokenAmount as 8-byte little-endian
	if tokenAmount.BitLen() > 64 {
		fmt.Println("Token amount exceeds 64 bits")
		return nil
	}
	// Convert to uint64 first since binary.Write handles native types
	tokenUint64 := tokenAmount.Uint64()
	if err := binary.Write(buf, binary.LittleEndian, tokenUint64); err != nil {
		fmt.Println("Error writing token amount:", err)
		return nil
	}

	// 3. Write sol_lamport_buy_amount as 8-byte little-endian
	if sol_lamport_buy_amount.BitLen() > 64 {
		fmt.Println("SOL amount exceeds 64 bits")
		return nil
	}

	sol_lamport_buy_amount = handlers.AddSlippageToBuy(sol_lamport_buy_amount, slippage)
	solUint64 := sol_lamport_buy_amount.Uint64()
	if err := binary.Write(buf, binary.LittleEndian, solUint64); err != nil {
		fmt.Println("Error writing SOL amount:", err)
		return nil
	}

	return buf.Bytes()
}
