// we need to be able to dictate what gets returned dependant on buy/sell
package instructions

import (
	"context"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/solana/instructions"
	solanaPda "personal_bot/internal/solana/pda"
	ammConstants "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm"
	pumpfunamm "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/instructions/pool"
	pumpfunAmmPda "personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/pda"
	"personal_bot/internal/solana/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type InformationFromAccounts struct {
	mintPoolAta           string
	wsolPoolAta           string
	poolData              pool.Pool
	associateTokenAddress string
}

func getBaseAccounts(ctx context.Context, token solana.PK, wallet solana.PK, httpClient *rpc.Client) (accounts []*solana.AccountMeta, err error, info InformationFromAccounts) {
	poolAddress, err := getPoolPDA(token)
	if err != nil {
		return
	}

	isNewToken, err := instructions.IsTokenAccountNew(ctx, token, httpClient)
	if err != nil {
		return
	}

	var baseMintATA solana.PublicKey
	if isNewToken {
		baseMintATA, _, err = solanaPda.FindToken2022AssociatedTokenAddress(wallet, token)
	} else {
		baseMintATA, _, err = solanaPda.FindTokenAssociatedTokenAddress(wallet, token)
	}
	if err != nil {
		return
	}

	wsolAddress := solana.MustPublicKeyFromBase58(constants.WSOLTokenAddress)
	userQuoteTokenAccount, _, err := solanaPda.FindTokenAssociatedTokenAddress(wallet, wsolAddress)
	if err != nil {
		return
	}

	poolBaseTokenAccount, err := getPoolTokenAccount(poolAddress, token.String(), isNewToken)
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

	poolBytes, err := pool.GetPoolDataBytes(ctx, poolAddressPK, httpClient)
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

	accounts = []*solana.AccountMeta{
		utils.GetAccountMeta(poolAddress, true, false),
		utils.GetAccountMeta(wallet.String(), true, true),
		utils.GetAccountMeta(pumpfunamm.GlobalConfig, false, false),
		utils.GetAccountMeta(token.String(), false, false),
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
	}
	return accounts, nil, InformationFromAccounts{poolData: poolData, wsolPoolAta: poolQuoteTokenAccount, mintPoolAta: poolBaseTokenAccount, associateTokenAddress: baseMintATA.String()}
}

func getPoolPDA(token solana.PK) (string, error) {
	return pumpfunAmmPda.GetPoolPda(token.String(), constants.WSOLTokenAddress)
}

func getPoolTokenAccount(poolAddress, mint string, isNewTokenAcc bool) (string, error) {
	return pumpfunAmmPda.GetPoolTokenAccount(poolAddress, mint, isNewTokenAcc)
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
