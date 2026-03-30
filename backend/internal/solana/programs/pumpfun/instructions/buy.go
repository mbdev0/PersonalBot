package instructions

import (
	"context"
	"fmt"
	"math/big"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/core/models"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/solana/instructions"
	bondingcurve "personal_bot/internal/solana/programs/pumpfun/bonding_curve"
	"personal_bot/internal/solana/programs/pumpfun/pda"
	"personal_bot/internal/solana/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

type AccountAddressesSet struct {
	BondingCurveAddress           string
	BondingCurveData              *models.BondingCurve
	CreatorAddress                string
	AssociatedBondingCurveAddress string
	AssociatedTokenAddressPubkey  solana.PublicKey
	TokenAddress                  string
	TokenProgram                  string
	WalletAddress                 string
	UserVolumeAccumulator         string
	BondingCurveV2Address         string
}

type BuyArgs struct {
	amount       uint64
	max_sol_cost uint64
	track_volume *bool
}

func GetBuyInstruction(buyTask *tasks.BuyTask, ctx context.Context) (instruction *solana.GenericInstruction, err error) {

	accountAddressesSet, err := setupAccountAddressSet(buyTask, ctx)
	if err != nil {
		return nil, fmt.Errorf("error setting up account address set: %w", err)
	}

	accounts := buildAccounts(accountAddressesSet)

	instructionData, err := createBuyData(*buyTask.BuyAmount, accountAddressesSet.BondingCurveData, buyTask.Slippage)
	if err != nil {
		return nil, fmt.Errorf("error creating buy data: %w", err)
	}

	progId := solana.MustPublicKeyFromBase58(constants.Program)
	buyInstructions := solana.NewInstruction(progId, accounts, instructionData)

	return buyInstructions, nil

}

func setupAccountAddressSet(buyTask *tasks.BuyTask, ctx context.Context) (AccountAddressesSet, error) {
	accountAddressesSet := &AccountAddressesSet{
		TokenAddress:  buyTask.Token.String(),
		WalletAddress: buyTask.Wallet.PublicKey().String(),
	}

	// Get and Set bonding curve information
	err := setBondingCurveInformation(accountAddressesSet, ctx, buyTask.HttpClient())
	if err != nil {
		return *accountAddressesSet, err
	}

	isNewTokenProgram, err := instructions.IsTokenAccountNew(buyTask.Token, ctx, buyTask.HttpClient())
	if err != nil {
		return *accountAddressesSet, err
	}

	if isNewTokenProgram {
		accountAddressesSet.TokenProgram = constants.Token2022Program
	} else {
		accountAddressesSet.TokenProgram = constants.TokenProgram
	}

	// Get and Set PDA account addresses
	err = resolvePDAs(accountAddressesSet, isNewTokenProgram)
	if err != nil {
		return *accountAddressesSet, err
	}

	return *accountAddressesSet, nil
}

func setBondingCurveInformation(accountAddressesSet *AccountAddressesSet, ctx context.Context, httpClient *rpc.Client) (err error) {
	bondingCurveAddress, err := pda.GetBondingCurveAddress(accountAddressesSet.TokenAddress)
	if err != nil {
		return fmt.Errorf("error getting bonding curve address: %w", err)
	}

	bondingCurveData, err, _ := bondingcurve.GetBondingCurveDataFromAddress(bondingCurveAddress, ctx, httpClient)
	if err != nil {
		return fmt.Errorf("error getting bonding curve data: %w", err)
	}

	accountAddressesSet.BondingCurveAddress = bondingCurveAddress
	accountAddressesSet.BondingCurveData = bondingCurveData

	return nil
}

func resolvePDAs(accountAddressesSet *AccountAddressesSet, isNewTokenAddress bool) (err error) {

	accountAddressesSet.CreatorAddress, err = pda.GetCreatorVaultAddress(accountAddressesSet.BondingCurveData.Creator.String())
	if err != nil {
		return fmt.Errorf("error getting creator vault address: %w", err)
	}

	accountAddressesSet.AssociatedBondingCurveAddress, err = pda.GetAssociatedBondingCurveAddress(accountAddressesSet.BondingCurveAddress, accountAddressesSet.TokenAddress, isNewTokenAddress)
	if err != nil {
		return fmt.Errorf("error getting associated bonding curve address: %w", err)
	}

	walletAddress := solana.MustPublicKeyFromBase58(accountAddressesSet.WalletAddress)
	tokenAddress := solana.MustPublicKeyFromBase58(accountAddressesSet.TokenAddress)

	if isNewTokenAddress {
		accountAddressesSet.AssociatedTokenAddressPubkey, _, err = pda.FindToken2022AssociatedTokenAddress(walletAddress, tokenAddress)
	} else {
		accountAddressesSet.AssociatedTokenAddressPubkey, _, err = pda.FindTokenAssociatedTokenAddress(walletAddress, tokenAddress)
	}
	if err != nil {
		return fmt.Errorf("error finding associated token address: %w", err)
	}

	accountAddressesSet.UserVolumeAccumulator, err = pda.GetUserVolumeAccumulatorAddress(walletAddress.String())
	if err != nil {
		return fmt.Errorf("error finding user volume accumululator address")
	}

	accountAddressesSet.BondingCurveV2Address, err = pda.GetBondingCurveV2Address(tokenAddress.String())
	if err != nil {
		return fmt.Errorf("error finding bonding-curve-v2 address")
	}

	return nil

}

func buildAccounts(accountAddressesSet AccountAddressesSet) (accounts []*solana.AccountMeta) {
	accounts = []*solana.AccountMeta{
		utils.GetAccountMeta(constants.GlobalAccount, true, false),
		utils.GetAccountMeta(constants.FeeRecipient, true, false),
		utils.GetAccountMeta(accountAddressesSet.TokenAddress, false, false),
		utils.GetAccountMeta(accountAddressesSet.BondingCurveAddress, true, false),
		utils.GetAccountMeta(accountAddressesSet.AssociatedBondingCurveAddress, true, false),
		utils.GetAccountMeta(accountAddressesSet.AssociatedTokenAddressPubkey.String(), true, false),
		utils.GetAccountMeta(accountAddressesSet.WalletAddress, true, true),
		utils.GetAccountMeta(solana.SystemProgramID.String(), false, false),
		utils.GetAccountMeta(accountAddressesSet.TokenProgram, false, false),
		utils.GetAccountMeta(accountAddressesSet.CreatorAddress, true, false),
		utils.GetAccountMeta(constants.EventAuthority, false, false),
		utils.GetAccountMeta(constants.Program, false, false),
		utils.GetAccountMeta(constants.GlobalVolumeAccumulator, true, false),
		utils.GetAccountMeta(accountAddressesSet.UserVolumeAccumulator, true, false),
		utils.GetAccountMeta(constants.FeeConfig, false, false),
		utils.GetAccountMeta(constants.FeeProgram, false, false),
		utils.GetAccountMeta(accountAddressesSet.BondingCurveV2Address, false, false),
	}
	return accounts
}

func createBuyData(solLamportBuyAmount big.Int, bondingCurveData *models.BondingCurve, slippage float64) (data []byte, err error) {
	tokenAmount, err, hasCompleted := bondingcurve.GetBuyTokenAmountFrom(solLamportBuyAmount, bondingCurveData)
	if err != nil || hasCompleted {
		if hasCompleted {
			return nil, fmt.Errorf("the coin has completed the bonding curve: %w", err)
		} else {
			return nil, fmt.Errorf("error getting buy token amount: %w", err)
		}
	}

	solLamportBuyAmount = addSlippageToBuy(solLamportBuyAmount, slippage)

	buyArgs := BuyArgs{
		amount:       tokenAmount.Uint64(),
		max_sol_cost: solLamportBuyAmount.Uint64(),
	}

	data, err = borsh.Serialize(buyArgs)
	if err != nil {
		return nil, err
	}

	data = append(constants.BuyInstructionDiscriminator[:], data...)

	return data, nil

}

func addSlippageToBuy(lamportAmount big.Int, slippagePercentage float64) (newBuyAmount big.Int) {
	slippageFloat := new(big.Float).SetFloat64(slippagePercentage)
	slippageFloat = new(big.Float).Add(slippageFloat, big.NewFloat(1))

	lamportFloat := new(big.Float).SetInt(&lamportAmount)

	newLamportAmountFloat := new(big.Float).Mul(slippageFloat, lamportFloat)
	newLamportAmountInt, _ := newLamportAmountFloat.Int(new(big.Int))

	return *newLamportAmountInt
}
