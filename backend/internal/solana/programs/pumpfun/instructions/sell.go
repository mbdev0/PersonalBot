package instructions

import (
	"context"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/core/models"
	"personal_bot/internal/core/position"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/solana/client"
	"personal_bot/internal/solana/instructions"
	bondingcurve "personal_bot/internal/solana/programs/pumpfun/bonding_curve"
	"personal_bot/internal/solana/programs/pumpfun/pda"
	"personal_bot/internal/solana/utils"
	"personal_bot/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

type SellArgs struct {
	amount         uint64
	min_sol_output uint64
}

var bondingCurveData *models.BondingCurve
var associatedTokenAddress solana.PublicKey

func GetSellInstruction(sellTask *tasks.SellTask, ctx context.Context, position *position.Position) (*solana.GenericInstruction, error) {
	accounts, err := getAccounts(sellTask, ctx)
	if err != nil {
		logger.Error("Error getting accounts for sell instruction", err)
		return nil, err
	}

	instructionData, err := getInstructionData(sellTask, ctx, position)

	if err != nil {
		logger.Error("Error getting instruction data for sell instruction", err)
		return nil, err
	}

	sellInstructions := solana.NewInstruction(solana.MustPublicKeyFromBase58(constants.Program), accounts, instructionData)
	return sellInstructions, nil
}

func getAccounts(sellTask *tasks.SellTask, ctx context.Context) ([]*solana.AccountMeta, error) {
	bondingCurveAddress, err := pda.GetBondingCurveAddress(sellTask.Token.String())
	if err != nil {
		logger.Error("Error getting bonding curve address:", err)
		return nil, err
	}

	isNewTokenAddress, err := instructions.IsTokenAccountNew(sellTask.Token, ctx, sellTask.HttpClient())
	if err != nil {
		return nil, err
	}

	var tokenProgram string
	if isNewTokenAddress {
		tokenProgram = constants.Token2022Program
	} else {
		tokenProgram = constants.TokenProgram
	}

	associatedBondingCurveAddress, err := pda.GetAssociatedBondingCurveAddress(bondingCurveAddress, sellTask.Token.String(), isNewTokenAddress)
	if err != nil {
		logger.Error("Error getting associated bonding curve address:", err)
		return nil, err
	}

	var ata solana.PublicKey
	if isNewTokenAddress {
		ata, _, err = pda.FindToken2022AssociatedTokenAddress(sellTask.Wallet.PublicKey(), sellTask.Token)
	} else {
		ata, _, err = pda.FindTokenAssociatedTokenAddress(sellTask.Wallet.PublicKey(), sellTask.Token)

	}

	if err != nil {
		logger.Error("Error getting token address: ", err)
		return nil, err
	}
	associatedTokenAddress = ata

	creatorAddress, err := getCreatorVaultAddress(bondingCurveAddress, ctx, sellTask.HttpClient())
	if err != nil {
		logger.Error("Error getting creator vault address:", err)
		return nil, err
	}

	bondingCurveV2Address, err := pda.GetBondingCurveV2Address(sellTask.Token.String())
	if err != nil {
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
		utils.GetAccountMeta(tokenProgram, false, false),
		utils.GetAccountMeta(constants.EventAuthority, false, false),
		utils.GetAccountMeta(constants.Program, false, false),
		utils.GetAccountMeta(constants.FeeConfig, false, false),
		utils.GetAccountMeta(constants.FeeProgram, false, false),
		utils.GetAccountMeta(bondingCurveV2Address, false, false),
	}

	return accounts, nil
}

func getCreatorVaultAddress(bondingCurveAddress string, ctx context.Context, httpClient *rpc.Client) (string, error) {
	data, err, _ := bondingcurve.GetBondingCurveDataFromAddress(bondingCurveAddress, ctx, httpClient)
	if err != nil {
		logger.Error("Error getting bonding curve data:", err)
		return "", err
	}
	bondingCurveData = data

	creatorAddress, err := pda.GetCreatorVaultAddress(bondingCurveData.Creator.String())

	if err != nil {
		logger.Error("Error getting creator address:", err)
		return "", err
	}

	return creatorAddress, nil
}

func getInstructionData(sellTask *tasks.SellTask, ctx context.Context, position *position.Position) ([]byte, error) {

	tokenAmount, solOutput, err := getTokenAmountAndSolOutput(sellTask, ctx, position)
	if err != nil {
		return nil, err
	}

	sellArgs := SellArgs{
		amount:         *tokenAmount,
		min_sol_output: *solOutput,
	}

	data, err := borsh.Serialize(sellArgs)
	if err != nil {
		return nil, err
	}

	return append(constants.SellInstructionDiscriminator[:], data...), nil
}

func getTokenAmountAndSolOutput(sellTask *tasks.SellTask, ctx context.Context, position *position.Position) (tokenAmount *uint64, solOutput *uint64, err error) {
	if position != nil {
		tokens, _ := position.TokenRemaining.Uint64()
		tokenAmount = &tokens
	} else {
		tokenAmount, err = client.GetTokenAccountBalance(associatedTokenAddress, sellTask.HttpClient(), ctx)
		if err != nil {
			return nil, nil, err
		}
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
