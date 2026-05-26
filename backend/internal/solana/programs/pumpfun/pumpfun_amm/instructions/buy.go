package instructions

import (
	"context"
	"math/big"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/core/tasks"
	ammConstants "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/constants"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/pool"

	"personal_bot/internal/solana/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type BuyArgs struct {
	SpendableQuoteIn uint64
	MinBaseAmountOut uint64
	TrackVolume      *bool
}

type accountsNeededForInstructions struct {
	mintPoolAta string
	wsolPoolAta string
}

func GetBuyInstruction(ctx context.Context, buyTask *tasks.BuyTask) (instruction *solana.GenericInstruction, err error) {
	accounts, err, instructionsForAccount := getAccounts(ctx, buyTask)
	if err != nil {
		return
	}

	instructions, err := getInstructionData(ctx, *buyTask, instructionsForAccount.mintPoolAta, instructionsForAccount.wsolPoolAta)
	if err != nil {
		return
	}

	progId := solana.MustPublicKeyFromBase58(constants.PumpFunAMMProgram)
	buyInstructions := solana.NewInstruction(progId, accounts, instructions)

	return buyInstructions, nil
}

// TODO - SELL uses mostly the same accounts - lets move it into a seperate file for all these pda getters
// TODO - Sell has everything except global and user volume accumalators
func getAccounts(ctx context.Context, task *tasks.BuyTask) (accounts []*solana.AccountMeta, err error, instructionAccounts accountsNeededForInstructions) {

	accounts, err, informationFromAccounts := getBaseAccounts(ctx, task.Token, task.Wallet.PublicKey(), task.HttpClient())
	if err != nil {
		return
	}

	userVolumeAccumulator, err := getUserVolumeAccumulator(task.Wallet.String())
	if err != nil {
		return
	}

	buyAccounts := []*solana.AccountMeta{
		utils.GetAccountMeta(ammConstants.GlobalVolumeAccumulator, false, false),
		utils.GetAccountMeta(userVolumeAccumulator, true, false),
		utils.GetAccountMeta(ammConstants.FeeConfig, false, false),
		utils.GetAccountMeta(constants.FeeProgram, false, false),
	}

	accounts = append(accounts, buyAccounts...)

	if informationFromAccounts.poolData.IsCashbackCoin {
		userVolumeAccumulatorWsolTokenAccount, err := getUserVolumeAccumulatorWsolTokenAccount(userVolumeAccumulator, constants.TokenProgram, constants.WSOLTokenAddress)
		if err != nil {
			return nil, err, accountsNeededForInstructions{}
		}
		accounts = append(accounts, utils.GetAccountMeta(userVolumeAccumulatorWsolTokenAccount, true, false))
	}

	poolV2, err := getPoolV2(task.GetToken())
	if err != nil {
		return
	}

	accounts = append(accounts, utils.GetAccountMeta(poolV2, false, false))
	accounts = append(accounts, utils.GetAccountMeta(ammConstants.BuyBackVault, false, false))
	accounts = append(accounts, utils.GetAccountMeta(ammConstants.BuyBackVaultWsol, true, false))

	return accounts, nil, accountsNeededForInstructions{mintPoolAta: informationFromAccounts.mintPoolAta, wsolPoolAta: informationFromAccounts.wsolPoolAta}
}

type poolBalances struct {
	tokenPoolBalance float64
	wsolPoolBalance  float64
}

func getInstructionData(ctx context.Context, task tasks.BuyTask, poolAtaAddress, wsolAtaAddress string) ([]byte, error) {
	poolBalances, err := pool.GetTokenBalances(ctx, poolAtaAddress, wsolAtaAddress, task.HttpClient())
	if err != nil {
		return nil, err
	}

	// convert pool balances to big.Float
	tokenPool := new(big.Float).SetFloat64(poolBalances.TokenPoolBalance)
	wsolPool := new(big.Float).SetFloat64(poolBalances.WsolPoolBalance)

	// k = tokenPool * wsolPool
	k := new(big.Float).Mul(tokenPool, wsolPool)

	// buyAmount in SOL (convert from lamports)
	lamportsConversion := new(big.Float).SetUint64(constants.LamportsConversion)
	buyAmountLamports := new(big.Float).SetUint64(task.BuyAmount.Uint64())
	buyAmountSol := new(big.Float).Quo(buyAmountLamports, lamportsConversion)

	// newWsolPool = wsolPool + buyAmountSol
	newWsolPool := new(big.Float).Add(wsolPool, buyAmountSol)

	// newTokenPool = k / newWsolPool
	newTokenPool := new(big.Float).Quo(k, newWsolPool)

	// tokensOut = tokenPool - newTokenPool
	tokensOut := new(big.Float).Sub(tokenPool, newTokenPool)

	// apply slippage: minBaseAmountOut = tokensOut * (1 - slippage) * tokenDecimals
	slippage := new(big.Float).SetFloat64(task.Slippage)
	one := new(big.Float).SetFloat64(1)
	slippageMultiplier := new(big.Float).Sub(one, slippage)

	tokenDecimals := new(big.Float).SetUint64(constants.TokenAmountDecimals)

	minBaseAmountOut := new(big.Float).Mul(tokensOut, slippageMultiplier)
	minBaseAmountOut.Mul(minBaseAmountOut, tokenDecimals)

	minBaseAmountOutUint64, _ := minBaseAmountOut.Uint64()

	buyArgs := BuyArgs{
		SpendableQuoteIn: task.BuyAmount.Uint64(),
		MinBaseAmountOut: minBaseAmountOutUint64,
	}

	serialised, err := borsh.Serialize(buyArgs)
	if err != nil {
		return nil, err
	}

	return append(ammConstants.BuyExactQuoteInDiscriminator, serialised...), nil
}

// func getTokenBalances(ctx context.Context, poolAtaAddress, wsolAtaAddress string, httpClient *rpc.Client) (poolBalances, error) {
// 	var tokenBalance, wsolBalance float64

// 	group, ctx := errgroup.WithContext(ctx)
// 	group.Go(
// 		func() error {
// 			poolAtaAddy, err := solana.PublicKeyFromBase58(poolAtaAddress)
// 			if err != nil {
// 				return err
// 			}

// 			res, err := httpClient.GetTokenAccountBalance(ctx, poolAtaAddy, rpc.CommitmentConfirmed)
// 			if err != nil {
// 				return err
// 			}

// 			poolAtaBalance, err := strconv.ParseFloat(res.Value.UiAmountString, 64)
// 			if err != nil {
// 				return err
// 			}

// 			tokenBalance = poolAtaBalance
// 			return nil
// 		},
// 	)

// 	group.Go(
// 		func() error {
// 			wsolAtaAddy, err := solana.PublicKeyFromBase58(wsolAtaAddress)
// 			if err != nil {
// 				return err
// 			}

// 			res, err := httpClient.GetTokenAccountBalance(ctx, wsolAtaAddy, rpc.CommitmentConfirmed)
// 			if err != nil {
// 				return err
// 			}

// 			wsolAtaBalance, err := strconv.ParseFloat(res.Value.UiAmountString, 32)
// 			if err != nil {
// 				return err
// 			}

// 			wsolBalance = wsolAtaBalance
// 			return nil
// 		},
// 	)

// 	if err := group.Wait(); err != nil {
// 		return poolBalances{}, err
// 	}

// 	return poolBalances{tokenPoolBalance: tokenBalance, wsolPoolBalance: wsolBalance}, nil
// }
