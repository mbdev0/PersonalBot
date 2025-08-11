package instructions

import (
	"bytes"
	"encoding/binary"
	"pump_fun/internal/core/constants"
	"pump_fun/internal/core/models"
	"pump_fun/internal/core/tasks"
	"pump_fun/internal/solana/client"
	bondingcurve "pump_fun/internal/solana/programs/pumpfun/bonding_curve"
	"pump_fun/internal/solana/programs/pumpfun/pda"
	"pump_fun/internal/solana/utils"
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go"
)

var bondingCurveData *models.BondingCurve
var associatedTokenAddress solana.PublicKey

func GetSellInstruction(sellTask *tasks.SellTask) (*solana.GenericInstruction, error) {
	accounts, err := getAccounts(sellTask)
	if err != nil {
		logger.Error("Error getting accounts for sell instruction", err)
		return nil, err
	}

	instructionData, err := getInstructionData(sellTask)

	if err != nil {
		logger.Error("Error getting instruction data for sell instruction", err)
		return nil, err
	}

	sell_instructions := solana.NewInstruction(solana.MustPublicKeyFromBase58(constants.Program), accounts, instructionData)
	return sell_instructions, nil
}

func getAccounts(sellTask *tasks.SellTask) ([]*solana.AccountMeta, error) {
	bondingCurveAddress, err := pda.GetBondingCurveAddress(sellTask.TokenAddress.String())
	if err != nil {
		logger.Error("Error getting bonding curve address:", err)
		return nil, err
	}

	associatedBondingCurveAddress, err := pda.GetAssociatedBondingCurveAddress(bondingCurveAddress, sellTask.TokenAddress.String())
	if err != nil {
		logger.Error("Error getting associated bonding curve address:", err)
		return nil, err
	}

	ATA, _, err := solana.FindAssociatedTokenAddress(sellTask.Wallet.PublicKey(), sellTask.TokenAddress)
	if err != nil {
		logger.Error("Error getting token address: ", err)
		return nil, err
	}
	associatedTokenAddress = ATA

	creatorAddress, err := getCreatorVaultAddress(bondingCurveAddress, sellTask.CancelToken)
	if err != nil {
		logger.Error("Error getting creator vault address:", err)
		return nil, err
	}

	accounts := []*solana.AccountMeta{
		utils.GetAccountMeta(constants.GlobalAccount, false, false),
		utils.GetAccountMeta(constants.FeeRecipient, true, false),
		utils.GetAccountMeta(sellTask.TokenAddress.String(), false, false),
		utils.GetAccountMeta(bondingCurveAddress, true, false),
		utils.GetAccountMeta(associatedBondingCurveAddress, true, false),
		utils.GetAccountMeta(associatedTokenAddress.String(), true, false),
		utils.GetAccountMeta(sellTask.Wallet.PublicKey().String(), true, true),
		utils.GetAccountMeta(solana.SystemProgramID.String(), false, false),
		utils.GetAccountMeta(creatorAddress, true, false),
		utils.GetAccountMeta(constants.TokenProgram, false, false),
		utils.GetAccountMeta(constants.EventAuthority, false, false),
		utils.GetAccountMeta(constants.Program, false, false),
	}

	return accounts, nil
}

func getCreatorVaultAddress(bondingCurveAddress string, cancellationToken models.CancelToken) (string, error) {
	data, err, _ := bondingcurve.GetBondingCurveDataFromAddress(bondingCurveAddress, cancellationToken)
	if err != nil {
		logger.Error("Error getting bonding curve data:", err)
		return "", err
	}
	bondingCurveData = data

	creatorAddress, err := pda.GetCreatorVaultAddress(bondingCurveData.DevWallet.String())

	if err != nil {
		logger.Error("Error getting creator address:", err)
		return "", err
	}

	return creatorAddress, nil
}

func getInstructionData(sellTask *tasks.SellTask) ([]byte, error) {

	tokenAmount, solOutput, err := getTokenAmountAndSolOutput(sellTask)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	discriminator := constants.SellInstructionDiscriminator
	if _, err := buf.Write(discriminator[:]); err != nil {
		logger.Error("Error writing discriminator to buffer", err)
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, tokenAmount); err != nil {
		logger.Error("Error writing token amount", err)
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, solOutput); err != nil {
		logger.Error("Error writing min_sol_output", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func getTokenAmountAndSolOutput(sellTask *tasks.SellTask) (tokenAmount *uint64, solOutput *uint64, err error) {
	tokenAmount, err = client.GetTokenAccountBalance(associatedTokenAddress, sellTask.CancelToken)
	if err != nil {
		return nil, nil, err
	}

	if sellTask.PercentageToSell > 0 && sellTask.PercentageToSell <= 1 {
		percentageToSell := sellTask.PercentageToSell
		*tokenAmount = uint64(float64(*tokenAmount) * percentageToSell)
	}

	sol_output := bondingcurve.GetSolanaTokenPrice(*bondingCurveData, *tokenAmount)
	slippage_sol_output := float64(*sol_output) * (1 - sellTask.Slippage)
	min_sol_output := uint64(slippage_sol_output)

	return tokenAmount, &min_sol_output, nil

}
