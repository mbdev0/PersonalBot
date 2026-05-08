package instructions

import (
	"context"
	"math/big"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/solana/instructions"
	solanaPda "personal_bot/internal/solana/pda"
	bondingcurve "personal_bot/internal/solana/programs/pumpfun/bonding_curve"
	"personal_bot/internal/solana/programs/pumpfun/pda"
	ammConstants "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm"
	pumpfunamm "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/instructions/pool"
	pumpfunAmmPda "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/pda"
	"strconv"

	"personal_bot/internal/solana/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"golang.org/x/sync/errgroup"
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

func getAccounts(ctx context.Context, task *tasks.BuyTask) (accounts []*solana.AccountMeta, err error, instructionAccounts accountsNeededForInstructions) {
	poolAddress, err := getPoolPDA(ctx, task)
	if err != nil {
		return
	}

	isNewToken, err := instructions.IsTokenAccountNew(ctx, solana.MustPublicKeyFromBase58(task.GetToken()), task.HttpClient())
	if err != nil {
		return
	}

	var baseMintATA solana.PublicKey
	if isNewToken {
		baseMintATA, _, err = solanaPda.FindToken2022AssociatedTokenAddress(task.Wallet.PublicKey(), task.Token)
	} else {
		baseMintATA, _, err = solanaPda.FindTokenAssociatedTokenAddress(task.Wallet.PublicKey(), task.Token)
	}
	if err != nil {
		return
	}

	// var userQuoteTokenAccount solana.PK
	wsolAddress := solana.MustPublicKeyFromBase58(constants.WSOLTokenAddress)
	// if isNewToken {
	// 	userQuoteTokenAccount, _, err = solanaPda.FindToken2022AssociatedTokenAddress(task.Wallet.PublicKey(), wsolAddress)
	// } else {}
	userQuoteTokenAccount, _, err := solanaPda.FindTokenAssociatedTokenAddress(task.Wallet.PublicKey(), wsolAddress)
	if err != nil {
		return
	}

	//8. poolBaseTokenAccount -> PDA(pool, base_token_program, base_mint)
	poolBaseTokenAccount, err := getPoolTokenAccount(poolAddress, task.GetToken(), isNewToken)
	if err != nil {
		return
	}

	poolQuoteTokenAccount, err := getPoolTokenAccount(poolAddress, constants.WSOLTokenAddress, false)
	if err != nil {
		return
	}

	protocolFeeRecipientTokenAccount, err := getProtocalFeeRecipientTokenAccount()
	if err != nil {
		return
	}
	poolAddressPK, err := solana.PublicKeyFromBase58(poolAddress)
	if err != nil {
		return
	}

	poolBytes, err := pool.GetPoolDataBytes(ctx, poolAddressPK, task.HttpClient())
	if err != nil {
		return
	}

	poolData, err := pool.DecodePoolData(poolBytes)
	if err != nil {
		return
	}

	creatorVaultAuthority, err := getCoinCreatorVaultAuthority(poolData.CoinCreator.String())
	if err != nil {
		return
	}

	coinVaultAta, err := getCoinCreatorVaultAta(creatorVaultAuthority, constants.WSOLTokenAddress)
	if err != nil {
		return
	}

	userVolumeAccumulator, err := getUserVolumeAccumulator(task.Wallet.PublicKey().String())
	if err != nil {
		return
	}

	accounts = []*solana.AccountMeta{
		utils.GetAccountMeta(poolAddress, true, false),
		utils.GetAccountMeta(task.GetWallet(), true, true),
		utils.GetAccountMeta(pumpfunamm.GlobalConfig, false, false),
		utils.GetAccountMeta(task.GetToken(), false, false),
		utils.GetAccountMeta(constants.WSOLTokenAddress, false, false),
		utils.GetAccountMeta(baseMintATA.String(), true, false),
		utils.GetAccountMeta(userQuoteTokenAccount.String(), true, false),
		utils.GetAccountMeta(poolBaseTokenAccount, true, false),
		utils.GetAccountMeta(poolQuoteTokenAccount, true, false),
		utils.GetAccountMeta(ammConstants.ProtocolFeeRecipient, false, false),
		utils.GetAccountMeta(protocolFeeRecipientTokenAccount, true, false),
		utils.GetAccountMeta(getTokenProgram(isNewToken), false, false),
		utils.GetAccountMeta(constants.TokenProgram, false, false),
		utils.GetAccountMeta(constants.SystemProgram, false, false),
		utils.GetAccountMeta(constants.AssociatedTokenProgram, false, false),
		utils.GetAccountMeta(ammConstants.EventAuthority, false, false),
		utils.GetAccountMeta(constants.PumpFunAMMProgram, false, false),
		utils.GetAccountMeta(coinVaultAta, true, false),
		utils.GetAccountMeta(creatorVaultAuthority, false, false),
		utils.GetAccountMeta(ammConstants.GlobalVolumeAccumulator, false, false),
		utils.GetAccountMeta(userVolumeAccumulator, true, false),
		utils.GetAccountMeta(ammConstants.FeeConfig, false, false),
		utils.GetAccountMeta(constants.FeeProgram, false, false),
	}

	if poolData.IsCashbackCoin {
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
	accounts = append(accounts, utils.GetAccountMeta("5cjcW9wExnJJiqgLjq7DEG75Pm6JBgE1hNv4B2vHXUW6", false, false))
	accounts = append(accounts, utils.GetAccountMeta("GYH1Gae1wJytMSvMvw8JVcv7nuAbxi8i9erNVbERnzXd", true, false))

	return accounts, nil, accountsNeededForInstructions{mintPoolAta: poolBaseTokenAccount, wsolPoolAta: poolQuoteTokenAccount}
}

func getPoolPDA(ctx context.Context, task *tasks.BuyTask) (string, error) {
	bondingCurveAddress, err := pda.GetBondingCurveAddress(task.GetToken(), constants.PumpFunProgram)
	if err != nil {
		return "", err
	}

	bondingCurveData, err, isComplete := bondingcurve.GetBondingCurveDataFromAddress(ctx, bondingCurveAddress, task.HttpClient())
	if err != nil && !isComplete {
		return "", err
	}

	return pumpfunAmmPda.GetPoolPda(bondingCurveData.Creator.String(), task.GetToken(), constants.WSOLTokenAddress)
}

func getPoolTokenAccount(poolAddress, mint string, isNewTokenAcc bool) (string, error) {
	var address string
	var err error
	if isNewTokenAcc {
		address, err = pumpfunAmmPda.GetPoolTokenAccount(poolAddress, mint, constants.Token2022Program)
	} else {
		address, err = pumpfunAmmPda.GetPoolTokenAccount(poolAddress, mint, constants.TokenProgram)
	}
	return address, err
}

func getProtocalFeeRecipientTokenAccount() (string, error) {
	address, err := pumpfunAmmPda.GetProtocalFeeRecipientTokenAccount(ammConstants.ProtocolFeeRecipient, constants.TokenProgram, constants.WSOLTokenAddress)

	return address, err
}

func getTokenProgram(isNewToken bool) string {
	if isNewToken {
		return constants.Token2022Program
	} else {
		return constants.TokenProgram
	}
}

func getCoinCreatorVaultAuthority(creator string) (string, error) {
	return pumpfunAmmPda.GetCoinCreatorVaultAuthority(creator)
}

func getCoinCreatorVaultAta(coinCreatorVaultAuthority, wsolMint string) (string, error) {
	return pumpfunAmmPda.GetCoinCreatorVaultAta(coinCreatorVaultAuthority, constants.TokenProgram, wsolMint)
}

func getUserVolumeAccumulator(wallet string) (string, error) {
	return pumpfunAmmPda.GetUserVolumeAccumulatorAddress(wallet)
}

func getUserVolumeAccumulatorWsolTokenAccount(userVolumeAccumulator, quoteTokenProgram, quoteMint string) (string, error) {
	return pumpfunAmmPda.GetUserVolumeAccumulatorWsolTokenAccount(userVolumeAccumulator, quoteTokenProgram, quoteMint)
}

func getPoolV2(baseMint string) (string, error) {
	return pumpfunAmmPda.GetPoolV2Account(baseMint)
}

func getBuybackFeeRecipientTokenAccount(quoteMint, buybackFeeRecipient, quoteTokenProgram string) (string, error) {
	return pumpfunAmmPda.GetBuybackFeeRecipientTokenAccount(quoteMint, buybackFeeRecipient, quoteTokenProgram)
}

type poolBalances struct {
	tokenPoolBalance float64
	wsolPoolBalance  float64
}

func getInstructionData(ctx context.Context, task tasks.BuyTask, poolAtaAddress, wsolAtaAddress string) ([]byte, error) {
	poolBalances, err := getTokenBalances(ctx, poolAtaAddress, wsolAtaAddress, task.HttpClient())
	if err != nil {
		return nil, err
	}

	// convert pool balances to big.Float
	tokenPool := new(big.Float).SetFloat64(poolBalances.tokenPoolBalance)
	wsolPool := new(big.Float).SetFloat64(poolBalances.wsolPoolBalance)

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
func getTokenBalances(ctx context.Context, poolAtaAddress, wsolAtaAddress string, httpClient *rpc.Client) (poolBalances, error) {
	var tokenBalance, wsolBalance float64

	group, ctx := errgroup.WithContext(ctx)
	group.Go(
		func() error {
			poolAtaAddy, err := solana.PublicKeyFromBase58(poolAtaAddress)
			if err != nil {
				return err
			}

			res, err := httpClient.GetTokenAccountBalance(ctx, poolAtaAddy, rpc.CommitmentConfirmed)
			if err != nil {
				return err
			}

			poolAtaBalance, err := strconv.ParseFloat(res.Value.UiAmountString, 64)
			if err != nil {
				return err
			}

			tokenBalance = poolAtaBalance
			return nil
		},
	)

	group.Go(
		func() error {
			wsolAtaAddy, err := solana.PublicKeyFromBase58(wsolAtaAddress)
			if err != nil {
				return err
			}

			res, err := httpClient.GetTokenAccountBalance(ctx, wsolAtaAddy, rpc.CommitmentConfirmed)
			if err != nil {
				return err
			}

			wsolAtaBalance, err := strconv.ParseFloat(res.Value.UiAmountString, 32)
			if err != nil {
				return err
			}

			wsolBalance = wsolAtaBalance
			return nil
		},
	)

	if err := group.Wait(); err != nil {
		return poolBalances{}, err
	}

	return poolBalances{tokenPoolBalance: tokenBalance, wsolPoolBalance: wsolBalance}, nil
}
