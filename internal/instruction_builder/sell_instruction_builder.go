package instructionbuilder

import (
	"bytes"
	"encoding/binary"
	"pump_fun/internal/logger"
	"pump_fun/internal/constants"
	"pump_fun/internal/monitoring/transactions/program_derived_address"
	"pump_fun/internal/monitoring/transactions/bonding_curve_decoder"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

func GetTransferInstruction(lamports uint64, publicKey solana.PublicKey, destinationAccount solana.PublicKey) *system.Instruction {
	instruction := system.NewTransferInstruction(
		lamports,
		publicKey,
		destinationAccount,
	).Build()

	return instruction
}

func GetSellInstruction(
	tokenAddress string,
	walletAddress string,
) (*solana.GenericInstruction, error) {	
	accounts, err := getAccounts(tokenAddress, walletAddress)
	if err != nil {
		logger.Error("Error getting accounts for sell instruction", err)
		return nil, err
	}

	instructionData, err := getInstructionData(26292925434, 582231)
	if err != nil {
		logger.Error("Error getting instruction data for sell instruction", err)
		return nil, err
	}
	
	sell_instructions := solana.NewInstruction(solana.MustPublicKeyFromBase58(constants.Program), accounts, instructionData)
	return sell_instructions, nil
}

func getAccounts(tokenAddress string, walletAddress string) ([]*solana.AccountMeta, error) {
	bondingCurveAddress, err := program_derived_address.GetBondingCurveAddress(tokenAddress)
	if err != nil {
		logger.Error("Error getting bonding curve address:", err)
		return nil, err
	}

	associatedBondingCurveAddress, err := program_derived_address.GetAssociatedBondingCurveAddress(bondingCurveAddress, tokenAddress)
	if err != nil {
		logger.Error("Error getting associated bonding curve address:", err)
		return nil, err
	}

	creatorAddress, err := getCreatorVaultAddress(bondingCurveAddress)
	if err != nil {
		logger.Error("Error getting creator vault address:", err)
		return nil, err
	}

	accounts := []*solana.AccountMeta{
		GetAccountMeta(constants.GlobalAccount, false, false),
		GetAccountMeta(constants.FeeRecipient, true, false),
		GetAccountMeta(tokenAddress, false, false),
		GetAccountMeta(bondingCurveAddress, true, false),
		GetAccountMeta(associatedBondingCurveAddress, true, false),
		GetAccountMeta("3VriQemVpwKTC1e4EbbmT7ZaNr9C7TfGbFmpWQgJxKyF", true, false), // This is the token account of the user
		GetAccountMeta(walletAddress, true, true),
		GetAccountMeta(solana.SystemProgramID.String(), false, false),
		GetAccountMeta(creatorAddress, true, false),
		GetAccountMeta(constants.TokenProgram, false, false),
		GetAccountMeta(constants.EventAuthority, false, false),
		GetAccountMeta(constants.Program, false, false),
	}

	return accounts, nil
}

func getCreatorVaultAddress(bondingCurveAddress string) (string, error) {
	bondingCurveData, err, _ := bonding_curve_decoder.GetBondingCurveDataFromAddress(bondingCurveAddress)
	if err != nil {
		logger.Error("Error getting bonding curve data:", err)
		return "", err
	}
	creatorAddress, err := program_derived_address.GetCreatorVaultAddress(bondingCurveData.DevWallet.String())
	if err != nil {
		logger.Error("Error getting creator address:", err)
		return "", err
	}

	return creatorAddress, nil
}

func getInstructionData(amount uint64, min_sol_output uint64) ([]byte, error) {
	buf := new(bytes.Buffer)

	discriminator := constants.SellInstructionDiscriminator
	if _, err := buf.Write(discriminator[:]); err != nil {
		logger.Error("Error writing discriminator to buffer", err)
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, amount); err != nil {
		logger.Error("Error writing token amount", err)
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, min_sol_output); err != nil {
		logger.Error("Error writing min_sol_output", err)
		return nil, err
	}

	return buf.Bytes(), nil
}