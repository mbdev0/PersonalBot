package instructions

import (
	"context"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/core/models"
	"personal_bot/internal/core/position"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/solana/client"
	"personal_bot/internal/solana/instructions"
	solanaPda "personal_bot/internal/solana/pda"
	bondingcurve "personal_bot/internal/solana/programs/pumpfun/bonding_curve"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_native/pda"

	"personal_bot/internal/solana/utils"
	"personal_bot/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type SellArgs struct {
	amount         uint64
	min_sol_output uint64
}

type BondingCurveInfo struct {
	bondingCurveAddress string
	bondingCurveData    *models.BondingCurve
}

type InstructionInfo struct {
	bondingCurveData       *models.BondingCurve
	associatedTokenAddress solana.PK
}

func GetSellInstruction(ctx context.Context, sellTask *tasks.SellTask, position *position.Position) (*solana.GenericInstruction, error) {
	bondingCurveAddress, err := pda.GetBondingCurveAddress(sellTask.Token.String())
	if err != nil {
		return nil, err
	}

	bondingCurveData, err, _ := bondingcurve.GetBondingCurveDataFromAddress(ctx, bondingCurveAddress, sellTask.HttpClient())
	if err != nil {
		return nil, err
	}

	accounts, instructionInfo, err := getAccounts(ctx, sellTask, BondingCurveInfo{bondingCurveAddress: bondingCurveAddress, bondingCurveData: bondingCurveData})
	if err != nil {
		return nil, err
	}

	instructionData, err := getInstructionData(ctx, sellTask, position, instructionInfo)
	if err != nil {
		return nil, err
	}
	sellInstructions := solana.NewInstruction(solana.MustPublicKeyFromBase58(constants.PumpFunProgram), accounts, instructionData)
	return sellInstructions, nil

}

func getAccounts(ctx context.Context, sellTask *tasks.SellTask, bondingCurveData BondingCurveInfo) (accounts []*solana.AccountMeta, instructionInfo InstructionInfo, err error) {

	tokenProgram, isNewTokenAddress, err := getTokenProgram(ctx, *sellTask)
	if err != nil {
		return
	}

	associatedBondingCurveAddress, err := GetAssociatedBondingCurveAddress(bondingCurveData.bondingCurveAddress, sellTask.GetToken(), isNewTokenAddress)
	if err != nil {
		return
	}

	creatorAddress, err := newGetCreatorVaultAddress(bondingCurveData.bondingCurveData)
	if err != nil {
		return
	}

	bondingCurveV2Address, err := getBondingCurveV2Address(sellTask.GetToken())
	if err != nil {
		return
	}

	associatedTokenAddress, err := getAssociatedTokenAddress(*sellTask, isNewTokenAddress)
	if err != nil {
		return
	}

	associatedTokenAddressPubKey, err := solana.PublicKeyFromBase58(associatedTokenAddress)
	if err != nil {
		return
	}

	feeRecipient := constants.FeeRecipient
	if bondingCurveData.bondingCurveData.IsMayhemMode == true {
		feeRecipient = constants.ReservedFeeRecipient
	}

	accounts = []*solana.AccountMeta{
		utils.GetAccountMeta(constants.GlobalAccount, false, false),
		utils.GetAccountMeta(feeRecipient, true, false),
		utils.GetAccountMeta(sellTask.Token.String(), false, false),
		utils.GetAccountMeta(bondingCurveData.bondingCurveAddress, true, false),
		utils.GetAccountMeta(associatedBondingCurveAddress, true, false),
		utils.GetAccountMeta(associatedTokenAddress, true, false),
		utils.GetAccountMeta(sellTask.Wallet.PublicKey().String(), true, true),
		utils.GetAccountMeta(solana.SystemProgramID.String(), false, false),
		utils.GetAccountMeta(creatorAddress, true, false),
		utils.GetAccountMeta(tokenProgram, false, false),
		utils.GetAccountMeta(constants.EventAuthority, false, false),
		utils.GetAccountMeta(constants.PumpFunProgram, false, false),
		utils.GetAccountMeta(constants.FeeConfig, false, false),
		utils.GetAccountMeta(constants.FeeProgram, false, false),
	}

	if bondingCurveData.bondingCurveData.IsCashbackCoin {
		userVolumeAccumulatorAddy, err := getUserVolumeAccumulatorAddress(sellTask.GetWallet())
		if err != nil {
			return nil, InstructionInfo{}, err
		}
		accounts = append(accounts, utils.GetAccountMeta(userVolumeAccumulatorAddy, true, false))
	}

	accounts = append(accounts, utils.GetAccountMeta(bondingCurveV2Address, false, false))
	accounts = append(accounts, utils.GetAccountMeta(constants.BuyBackFeeRecipient, true, false))

	return accounts, InstructionInfo{associatedTokenAddress: associatedTokenAddressPubKey, bondingCurveData: bondingCurveData.bondingCurveData}, nil
}

func getTokenProgram(ctx context.Context, sellTask tasks.SellTask) (tokenProgram string, isNewTokenAddress bool, err error) {
	isNewTokenAddress, err = instructions.IsTokenAccountNew(ctx, sellTask.Token, sellTask.HttpClient())
	if err != nil {
		return
	}

	if isNewTokenAddress {
		tokenProgram = constants.Token2022Program
	} else {
		tokenProgram = constants.TokenProgram
	}

	return tokenProgram, isNewTokenAddress, nil
}

func GetAssociatedBondingCurveAddress(bondingCurveAddress string, token string, isNewTokenAddress bool) (string, error) {
	associatedBondingCurveAddress, err := pda.GetAssociatedBondingCurveAddress(bondingCurveAddress, token, isNewTokenAddress)
	if err != nil {
		logger.Error("Error getting associated bonding curve address:", err)
		return "", err
	}

	return associatedBondingCurveAddress, nil

}

func newGetCreatorVaultAddress(bondingCurve *models.BondingCurve) (string, error) {
	creatorAddress, err := pda.GetCreatorVaultAddress(bondingCurve.Creator.String())

	if err != nil {
		logger.Error("Error getting creator address:", err)
		return "", err
	}

	return creatorAddress, nil
}

func getBondingCurveV2Address(token string) (string, error) {
	bondingCurveV2Address, err := pda.GetBondingCurveV2Address(token)
	if err != nil {
		return "", err
	}

	return bondingCurveV2Address, nil
}

func getUserVolumeAccumulatorAddress(wallet string) (string, error) {
	userVolumeAccumulatorAddy, err := pda.GetUserVolumeAccumulatorAddress(wallet)
	if err != nil {
		return "", err
	}

	return userVolumeAccumulatorAddy, nil
}

func getAssociatedTokenAddress(sellTask tasks.SellTask, isNewTokenProgram bool) (ata string, err error) {
	if isNewTokenProgram {
		ataPubKey, _, err := solanaPda.FindToken2022AssociatedTokenAddress(sellTask.Wallet.PublicKey(), sellTask.Token)
		if err != nil {
			return "", err
		}
		return ataPubKey.String(), nil
	} else {
		ataPubKey, _, err := solanaPda.FindTokenAssociatedTokenAddress(sellTask.Wallet.PublicKey(), sellTask.Token)
		if err != nil {
			return "", err
		}
		return ataPubKey.String(), nil

	}
}

func getInstructionData(ctx context.Context, sellTask *tasks.SellTask, position *position.Position, instructionInfo InstructionInfo) ([]byte, error) {
	tokenAmount, solOutput, err := getTokenAmountAndSolOutput(ctx, sellTask, position, instructionInfo)
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

func getTokenAmountAndSolOutput(ctx context.Context, sellTask *tasks.SellTask, position *position.Position, instructionInfo InstructionInfo) (tokenAmount *uint64, solOutput *uint64, err error) {
	if position != nil {
		tokens, _ := position.TokenRemaining.Uint64()
		tokenAmount = &tokens
	} else {
		tokenAmount, err = client.GetTokenAccountBalance(ctx, instructionInfo.associatedTokenAddress, sellTask.HttpClient())
		if err != nil {
			return nil, nil, err
		}
	}

	if sellTask.SellPercentage > 0 && sellTask.SellPercentage <= 1 {
		percentageToSell := sellTask.SellPercentage
		*tokenAmount = uint64(float64(*tokenAmount) * percentageToSell)
	}

	solAmnt := bondingcurve.GetSolanaTokenPrice(*instructionInfo.bondingCurveData, *tokenAmount)
	slippageSolOutput := float64(*solAmnt) * (1 - sellTask.Slippage)
	minSolOutput := uint64(slippageSolOutput)

	return tokenAmount, &minSolOutput, nil

}
