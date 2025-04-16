package instructionbuilder

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/big"
	"pump_fun/internal/constants"
	"pump_fun/internal/handlers"
	"pump_fun/internal/logger"
	"pump_fun/internal/monitoring/transactions/bonding_curve_decoder"

	"github.com/gagliardetto/solana-go"
)

func GetBuyInstruction(tokenAddress string, walletAddress string, sol_lamport_buy_amount big.Int, slippage float64) (instruction *solana.GenericInstruction, err error) {

	bondingCurveAddress, err := bonding_curve_decoder.GetBondingCurveAddress(tokenAddress)
	if err != nil {
		return nil, err
	}

	associatedBondingCurveAddress, err := bonding_curve_decoder.GetAssociatedBondingCurveAddress(bondingCurveAddress, tokenAddress)
	if err != nil {
		return nil, err
	}

	associatedTokenAddressPubkey, _, err := solana.FindAssociatedTokenAddress(solana.MustPublicKeyFromBase58(walletAddress), solana.MustPublicKeyFromBase58(tokenAddress))
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
		GetAccountMeta(constants.RentProgram, false, false),
		GetAccountMeta(constants.EventAuthority, false, false),
		GetAccountMeta(constants.Program, false, false),
	}

	instruction_data, err := create_buy_data(sol_lamport_buy_amount, bondingCurveAddress, slippage)
	if err != nil {
		logger.Log(logger.LevelError, "Error creating buy data", logger.String("error", err.Error()))
		return nil, err
	}

	buy_instructions := solana.NewInstruction(solana.MustPublicKeyFromBase58(constants.Program), accounts, instruction_data)

	return buy_instructions, nil

}

func create_buy_data(sol_lamport_buy_amount big.Int, bondingCurveAddr string, slippage float64) (data []byte, err error) {
	// Get the token amount
	tokenAmount, err, hasCompleted := handlers.GetBuyTokenAmountFrom(sol_lamport_buy_amount, bondingCurveAddr)
	if err != nil || hasCompleted {
		if hasCompleted {
			logger.Log(logger.LevelError, "The coin has completed the bonding curve")
			return nil, errors.New("the coin has completed the bonding curve")
		} else {
			logger.Log(logger.LevelError, "There was an error getting the token amount", logger.String("error", err.Error()))
			return nil, err
		}
	}

	// Create a buffer with capacity for:
	// 8 bytes discriminator + 8 bytes tokenAmount + 8 bytes sol_lamport_buy_amount
	buf := new(bytes.Buffer)

	discriminator := constants.BuyInstructionDiscriminator
	if _, err := buf.Write(discriminator[:]); err != nil {
		logger.Log(logger.LevelError, "Error writing discriminator", logger.String("error", err.Error()))
		return nil, err
	}

	if tokenAmount.BitLen() > 64 {
		logger.Log(logger.LevelError, "Token amount exceeds 64 bits")
		return nil, err
	}
	tokenUint64 := tokenAmount.Uint64()
	if err := binary.Write(buf, binary.LittleEndian, tokenUint64); err != nil {
		logger.Log(logger.LevelError, "Error writing token amount", logger.String("error", err.Error()))
		return nil, err
	}

	if sol_lamport_buy_amount.BitLen() > 64 {
		logger.Log(logger.LevelError, "SOL amount exceeds 64 bits")
		return nil, err
	}

	sol_lamport_buy_amount = handlers.AddSlippageToBuy(sol_lamport_buy_amount, slippage)
	solUint64 := sol_lamport_buy_amount.Uint64()
	if err := binary.Write(buf, binary.LittleEndian, solUint64); err != nil {
		logger.Log(logger.LevelError, "Error writing SOL amount", logger.String("error", err.Error()))
		return nil, err
	}

	return buf.Bytes(), err
}
