package instructions

import (
	"context"
	"math/big"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/core/position"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/solana/client"
	ammConstants "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/constants"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/pool"
	"personal_bot/internal/solana/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type SellArgs struct {
	BaseAmountIn      uint64
	MinQuoteAmountOut uint64
}

type sellInstructionInfo struct {
	mintPoolAta            string
	wsolPoolAta            string
	associatedTokenAddress string
}

func GetSellInstruction(ctx context.Context, sellTask *tasks.SellTask, pos *position.Position) (*solana.GenericInstruction, error) {
	accounts, information, err := getSellAccounts(ctx, sellTask)
	if err != nil {
		return nil, err
	}

	instructionData, err := getSellInstructionData(ctx, sellTask, pos, sellInstructionInfo{mintPoolAta: information.mintPoolAta, wsolPoolAta: information.wsolPoolAta})
	if err != nil {
		return nil, err
	}

	ammProgram := solana.MustPublicKeyFromBase58(constants.PumpFunAMMProgram)
	inst := solana.NewInstruction(ammProgram, accounts, instructionData)

	return inst, nil

}

func getSellAccounts(ctx context.Context, sellTask *tasks.SellTask) ([]*solana.AccountMeta, InformationFromAccounts, error) {
	baseAccounts, err, informationFromAccounts := getBaseAccounts(ctx, sellTask.Token, sellTask.Wallet.PublicKey(), sellTask.HttpClient())
	if err != nil {
		return nil, informationFromAccounts, err
	}

	baseAccounts = append(baseAccounts,
		utils.GetAccountMeta(ammConstants.FeeConfig, false, false),
		utils.GetAccountMeta(constants.FeeProgram, false, false),
	)

	poolData := informationFromAccounts.poolData

	remainingAccounts := []*solana.AccountMeta{}

	if poolData.IsCashbackCoin {
		userVolumeAccumulator, err := getUserVolumeAccumulator(sellTask.Wallet.PublicKey().String())
		if err != nil {
			return nil, informationFromAccounts, err
		}

		userVolumeAccumulatorWsolAta, err := getUserVolumeAccumulatorWsolTokenAccount(
			userVolumeAccumulator,
			constants.TokenProgram,
			constants.WSOLTokenAddress,
		)
		if err != nil {
			return nil, informationFromAccounts, err
		}

		remainingAccounts = append(remainingAccounts,
			utils.GetAccountMeta(userVolumeAccumulatorWsolAta, true, false),
			utils.GetAccountMeta(userVolumeAccumulator, true, false),
		)
	}

	if !poolData.CoinCreator.IsZero() {
		poolV2, err := getPoolV2(sellTask.Token.String())
		if err != nil {
			return nil, informationFromAccounts, err
		}
		remainingAccounts = append(remainingAccounts,
			utils.GetAccountMeta(poolV2, false, false),
		)
	}

	remainingAccounts = append(remainingAccounts,
		utils.GetAccountMeta(ammConstants.BuyBackVault, false, false),
		utils.GetAccountMeta(ammConstants.BuyBackVaultWsol, true, false),
	)

	baseAccounts = append(baseAccounts, remainingAccounts...)

	return baseAccounts, informationFromAccounts, nil
}

func getSellInstructionData(ctx context.Context, sellTask *tasks.SellTask, position *position.Position, instructionInfo sellInstructionInfo) ([]byte, error) {
	tokenAmount, solOutput, err := getTokenAmountAndSolOutput(ctx, sellTask, position, instructionInfo)
	if err != nil {
		return nil, err
	}

	sellArgs := SellArgs{
		BaseAmountIn:      tokenAmount,
		MinQuoteAmountOut: solOutput,
	}

	data, err := borsh.Serialize(sellArgs)
	if err != nil {
		return nil, err
	}

	return append(constants.SellInstructionDiscriminator[:], data...), nil
}

func getTokenAmountAndSolOutput(ctx context.Context, sellTask *tasks.SellTask, position *position.Position, instructionInfo sellInstructionInfo) (tokenAmount uint64, solOutput uint64, err error) {
	// resolve raw token amount
	var rawTokenAmount uint64
	if position != nil {
		rawTokenAmount, _ = position.TokenRemaining.Uint64()
	} else {
		ata, err := solana.PublicKeyFromBase58(instructionInfo.associatedTokenAddress)
		if err != nil {
			return 0, 0, err
		}

		balance, err := client.GetTokenAccountBalance(ctx, ata, sellTask.HttpClient())
		if err != nil {
			return 0, 0, err
		}
		rawTokenAmount = *balance
	}

	// apply sell percentage
	if sellTask.SellPercentage > 0 && sellTask.SellPercentage <= 1 {
		rawTokenAmount = uint64(float64(rawTokenAmount) * sellTask.SellPercentage)
	}

	// fetch pool balances
	poolBalances, err := pool.GetTokenBalances(ctx, instructionInfo.mintPoolAta, instructionInfo.wsolPoolAta, sellTask.HttpClient())
	if err != nil {
		return 0, 0, err
	}

	// constant product formula — reverse of buy
	// tokensIn → solOut
	tokenPool := new(big.Float).SetFloat64(poolBalances.TokenPoolBalance)
	wsolPool := new(big.Float).SetFloat64(poolBalances.WsolPoolBalance)

	k := new(big.Float).Mul(tokenPool, wsolPool)

	// convert token amount from raw to UI amount
	tokenDecimals := new(big.Float).SetUint64(constants.TokenAmountDecimals)
	tokensIn := new(big.Float).Quo(new(big.Float).SetUint64(rawTokenAmount), tokenDecimals)

	// newTokenPool = tokenPool + tokensIn
	newTokenPool := new(big.Float).Add(tokenPool, tokensIn)

	// newWsolPool = k / newTokenPool
	newWsolPool := new(big.Float).Quo(k, newTokenPool)

	// solOut = wsolPool - newWsolPool
	solOut := new(big.Float).Sub(wsolPool, newWsolPool)

	// convert to lamports and apply slippage
	lamportsConversion := new(big.Float).SetUint64(constants.LamportsConversion)
	solOutLamports := new(big.Float).Mul(solOut, lamportsConversion)

	slippage := new(big.Float).SetFloat64(sellTask.Slippage)
	one := new(big.Float).SetFloat64(1)
	slippageMultiplier := new(big.Float).Sub(one, slippage)

	minSolOutput := new(big.Float).Mul(solOutLamports, slippageMultiplier)
	minSolOutputUint64, _ := minSolOutput.Uint64()

	return rawTokenAmount, minSolOutputUint64, nil
}
