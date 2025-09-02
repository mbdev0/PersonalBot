package instructions

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"pump_fun/internal/core/constants"
	"pump_fun/internal/core/models"
	"pump_fun/internal/core/tasks"
	bondingcurve "pump_fun/internal/solana/programs/pumpfun/bonding_curve"
	"pump_fun/internal/solana/programs/pumpfun/pda"
	"pump_fun/internal/solana/utils"

	"github.com/gagliardetto/solana-go"
)

type AccountAddressesSet struct {
	BondingCurveAddress           string
	BondingCurveData              *models.BondingCurve
	CreatorAddress                string
	AssociatedBondingCurveAddress string
	AssociatedTokenAddressPubkey  solana.PublicKey
	TokenAddress                  string
	WalletAddress                 string
	UserVolumeAccumulator         string
}

func GetBuyInstruction(buyTask *tasks.BuyTask) (instruction *solana.GenericInstruction, err error) {

	accountAddressesSet, err := setupAccountAddressSet(buyTask)
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

func setupAccountAddressSet(buyTask *tasks.BuyTask) (AccountAddressesSet, error) {
	accountAddressesSet := &AccountAddressesSet{
		TokenAddress:  buyTask.Token.String(),
		WalletAddress: buyTask.Wallet.PublicKey().String(),
	}

	// Get and Set bonding curve information
	err := setBondingCurveInformation(accountAddressesSet, buyTask.Ctx())
	if err != nil {
		return *accountAddressesSet, err
	}

	// Get and Set PDA account addresses
	err = resolvePDAs(accountAddressesSet)
	if err != nil {
		return *accountAddressesSet, err
	}

	return *accountAddressesSet, nil
}

func setBondingCurveInformation(accountAddressesSet *AccountAddressesSet, ctx context.Context) (err error) {
	bondingCurveAddress, err := pda.GetBondingCurveAddress(accountAddressesSet.TokenAddress)
	if err != nil {
		return fmt.Errorf("error getting bonding curve address: %w", err)
	}

	bondingCurveData, err, _ := bondingcurve.GetBondingCurveDataFromAddress(bondingCurveAddress, ctx)
	if err != nil {
		return fmt.Errorf("error getting bonding curve data: %w", err)
	}

	accountAddressesSet.BondingCurveAddress = bondingCurveAddress
	accountAddressesSet.BondingCurveData = bondingCurveData

	return nil
}

func resolvePDAs(accountAddressesSet *AccountAddressesSet) (err error) {

	accountAddressesSet.CreatorAddress, err = pda.GetCreatorVaultAddress(accountAddressesSet.BondingCurveData.DevWallet.String())
	if err != nil {
		return fmt.Errorf("error getting creator vault address: %w", err)
	}

	accountAddressesSet.AssociatedBondingCurveAddress, err = pda.GetAssociatedBondingCurveAddress(accountAddressesSet.BondingCurveAddress, accountAddressesSet.TokenAddress)
	if err != nil {
		return fmt.Errorf("error getting associated bonding curve address: %w", err)
	}

	walletAddress := solana.MustPublicKeyFromBase58(accountAddressesSet.WalletAddress)
	tokenAddress := solana.MustPublicKeyFromBase58(accountAddressesSet.TokenAddress)

	accountAddressesSet.AssociatedTokenAddressPubkey, _, err = solana.FindAssociatedTokenAddress(walletAddress, tokenAddress)
	if err != nil {
		return fmt.Errorf("error finding associated token address: %w", err)
	}

	accountAddressesSet.UserVolumeAccumulator, err = pda.GetUserVolumeAccumulatorAddress(walletAddress.String())

	if err != nil {
		return fmt.Errorf("error finding user volume accumululator address")
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
		utils.GetAccountMeta(constants.TokenProgram, false, false),
		utils.GetAccountMeta(accountAddressesSet.CreatorAddress, true, false),
		utils.GetAccountMeta(constants.EventAuthority, false, false),
		utils.GetAccountMeta(constants.Program, false, false),
		utils.GetAccountMeta(constants.GlobalVolumeAccumulator, true, false),
		utils.GetAccountMeta(accountAddressesSet.UserVolumeAccumulator, true, false),
		utils.GetAccountMeta(constants.FeeConfig, false, false),
		utils.GetAccountMeta(constants.FeeProgram, false, false),
	}
	return accounts
}

func createBuyData(solLamportBuyAmount big.Int, bondingCurveData *models.BondingCurve, slippage float64) (data []byte, err error) {
	// Get the token amount
	tokenAmount, err, hasCompleted := bondingcurve.GetBuyTokenAmountFrom(solLamportBuyAmount, bondingCurveData)
	if err != nil || hasCompleted {
		if hasCompleted {
			return nil, fmt.Errorf("the coin has completed the bonding curve: %w", err)
		} else {
			return nil, fmt.Errorf("error getting buy token amount: %w", err)
		}
	}

	// Create a buffer with capacity for:
	// 8 bytes discriminator + 8 bytes tokenAmount + 8 bytes sol_lamport_buy_amount
	buf := new(bytes.Buffer)

	discriminator := constants.BuyInstructionDiscriminator
	if _, err := buf.Write(discriminator[:]); err != nil {
		return nil, fmt.Errorf("error writing discriminator to buffer: %w", err)
	}

	if tokenAmount.BitLen() > 64 {
		return nil, fmt.Errorf("token amount %d overflows 64 bits", tokenAmount)
	}
	tokenUint64 := tokenAmount.Uint64()
	if err := binary.Write(buf, binary.LittleEndian, tokenUint64); err != nil {
		return nil, fmt.Errorf("error writing token amount to buffer: %w", err)
	}

	if solLamportBuyAmount.BitLen() > 64 {
		return nil, fmt.Errorf("sol lamport buy amount exceeds 64 bits: %s", solLamportBuyAmount.String())
	}

	solLamportBuyAmount = addSlippageToBuy(solLamportBuyAmount, slippage)
	solUint64 := solLamportBuyAmount.Uint64()
	if err := binary.Write(buf, binary.LittleEndian, solUint64); err != nil {
		return nil, fmt.Errorf("error writing sol lamport buy amount to buffer: %w", err)
	}

	return buf.Bytes(), nil
}

func addSlippageToBuy(lamportAmount big.Int, slippagePercentage float64) (newBuyAmount big.Int) {
	slippageFloat := new(big.Float).SetFloat64(slippagePercentage)
	slippageFloat = new(big.Float).Add(slippageFloat, big.NewFloat(1))

	lamportFloat := new(big.Float).SetInt(&lamportAmount)

	newLamportAmountFloat := new(big.Float).Mul(slippageFloat, lamportFloat)
	newLamportAmountInt, _ := newLamportAmountFloat.Int(new(big.Int))

	return *newLamportAmountInt
}
