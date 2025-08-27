package instructions

import (
	"bytes"
	"context"
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

	sellInstructions := solana.NewInstruction(solana.MustPublicKeyFromBase58(constants.Program), accounts, instructionData)
	return sellInstructions, nil
}

func getAccounts(sellTask *tasks.SellTask) ([]*solana.AccountMeta, error) {
	bondingCurveAddress, err := pda.GetBondingCurveAddress(sellTask.Token.String())
	if err != nil {
		logger.Error("Error getting bonding curve address:", err)
		return nil, err
	}

	associatedBondingCurveAddress, err := pda.GetAssociatedBondingCurveAddress(bondingCurveAddress, sellTask.Token.String())
	if err != nil {
		logger.Error("Error getting associated bonding curve address:", err)
		return nil, err
	}

	ATA, _, err := solana.FindAssociatedTokenAddress(sellTask.Wallet.PublicKey(), sellTask.Token)
	if err != nil {
		logger.Error("Error getting token address: ", err)
		return nil, err
	}
	associatedTokenAddress = ATA

	creatorAddress, err := getCreatorVaultAddress(bondingCurveAddress, sellTask.Ctx())
	if err != nil {
		logger.Error("Error getting creator vault address:", err)
		return nil, err
	}

	accounts := []*solana.AccountMeta{
		utils.GetAccountMeta(constants.GlobalAccount, false, false),
		utils.GetAccountMeta(constants.FeeRecipient, true, false),
		utils.GetAccountMeta(sellTask.Token.String(), false, false),
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

func getCreatorVaultAddress(bondingCurveAddress string, ctx context.Context) (string, error) {
	data, err, _ := bondingcurve.GetBondingCurveDataFromAddress(bondingCurveAddress, ctx)
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
	tokenAmount, err = client.GetTokenAccountBalance(associatedTokenAddress, sellTask.Ctx())
	if err != nil {
		return nil, nil, err
	}

	if sellTask.SellPercentage > 0 && sellTask.SellPercentage <= 1 {
		percentageToSell := sellTask.SellPercentage
		*tokenAmount = uint64(float64(*tokenAmount) * percentageToSell)
	}

	solAmnt := bondingcurve.GetSolanaTokenPrice(*bondingCurveData, *tokenAmount)
	slippageSolOutput := float64(*solAmnt) * (1 - sellTask.Slippage)
	minSolOutput := uint64(slippageSolOutput)

	return tokenAmount, &minSolOutput, nil

}
