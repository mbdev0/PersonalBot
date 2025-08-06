package instructionbuilder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"pump_fun/internal/constants"
	"pump_fun/internal/handlers"
	"pump_fun/internal/models"
	"pump_fun/internal/models/tasks"
	"pump_fun/internal/monitoring/transactions/bonding_curve_decoder"
	"pump_fun/internal/monitoring/transactions/program_derived_address"

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

	instructionData, err := createBuyData(buyTask.BuyAmount, accountAddressesSet.BondingCurveData, buyTask.Slippage)
	if err != nil {
		return nil, fmt.Errorf("error creating buy data: %w", err)
	}

	progId := solana.MustPublicKeyFromBase58(constants.Program)
	buyInstructions := solana.NewInstruction(progId, accounts, instructionData)

	return buyInstructions, nil

}

func setupAccountAddressSet(buyTask *tasks.BuyTask) (AccountAddressesSet, error) {
	accountAddressesSet := &AccountAddressesSet{
		TokenAddress:  buyTask.TokenAddress.String(),
		WalletAddress: buyTask.Wallet.PublicKey().String(),
	}

	// Get and Set bonding curve information
	err := setBondingCurveInformation(accountAddressesSet, buyTask.CancelToken)
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

func setBondingCurveInformation(accountAddressesSet *AccountAddressesSet, cancellationToken models.CancelToken) (err error) {
	bondingCurveAddress, err := program_derived_address.GetBondingCurveAddress(accountAddressesSet.TokenAddress)
	if err != nil {
		return fmt.Errorf("error getting bonding curve address: %w", err)
	}

	bondingCurveData, err, _ := bonding_curve_decoder.GetBondingCurveDataFromAddress(bondingCurveAddress, cancellationToken)
	if err != nil {
		return fmt.Errorf("error getting bonding curve data: %w", err)
	}

	accountAddressesSet.BondingCurveAddress = bondingCurveAddress
	accountAddressesSet.BondingCurveData = bondingCurveData

	return nil
}

func resolvePDAs(accountAddressesSet *AccountAddressesSet) (err error) {

	accountAddressesSet.CreatorAddress, err = program_derived_address.GetCreatorVaultAddress(accountAddressesSet.BondingCurveData.DevWallet.String())
	if err != nil {
		return fmt.Errorf("error getting creator vault address: %w", err)
	}

	accountAddressesSet.AssociatedBondingCurveAddress, err = program_derived_address.GetAssociatedBondingCurveAddress(accountAddressesSet.BondingCurveAddress, accountAddressesSet.TokenAddress)
	if err != nil {
		return fmt.Errorf("error getting associated bonding curve address: %w", err)
	}

	walletAddress := solana.MustPublicKeyFromBase58(accountAddressesSet.WalletAddress)
	tokenAddress := solana.MustPublicKeyFromBase58(accountAddressesSet.TokenAddress)

	accountAddressesSet.AssociatedTokenAddressPubkey, _, err = solana.FindAssociatedTokenAddress(walletAddress, tokenAddress)
	if err != nil {
		return fmt.Errorf("error finding associated token address: %w", err)
	}

	accountAddressesSet.UserVolumeAccumulator, err = program_derived_address.GetUserVolumeAccumulatorAddress(walletAddress.String())

	if err != nil {
		return fmt.Errorf("error finding user volume accumululator address")
	}

	return nil

}

func buildAccounts(accountAddressesSet AccountAddressesSet) (accounts []*solana.AccountMeta) {
	accounts = []*solana.AccountMeta{
		GetAccountMeta(constants.GlobalAccount, true, false),
		GetAccountMeta(constants.FeeRecipient, true, false),
		GetAccountMeta(accountAddressesSet.TokenAddress, false, false),
		GetAccountMeta(accountAddressesSet.BondingCurveAddress, true, false),
		GetAccountMeta(accountAddressesSet.AssociatedBondingCurveAddress, true, false),
		GetAccountMeta(accountAddressesSet.AssociatedTokenAddressPubkey.String(), true, false),
		GetAccountMeta(accountAddressesSet.WalletAddress, true, true),
		GetAccountMeta(solana.SystemProgramID.String(), false, false),
		GetAccountMeta(constants.TokenProgram, false, false),
		GetAccountMeta(accountAddressesSet.CreatorAddress, true, false),
		GetAccountMeta(constants.EventAuthority, false, false),
		GetAccountMeta(constants.Program, false, false),
		GetAccountMeta(constants.GlobalVolumeAccumulator, true, false),
		GetAccountMeta(accountAddressesSet.UserVolumeAccumulator, true, false),
	}
	return accounts
}

func createBuyData(solLamportBuyAmount big.Int, bondingCurveData *models.BondingCurve, slippage float64) (data []byte, err error) {
	// Get the token amount
	tokenAmount, err, hasCompleted := handlers.GetBuyTokenAmountFrom(solLamportBuyAmount, bondingCurveData)
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

	solLamportBuyAmount = handlers.AddSlippageToBuy(solLamportBuyAmount, slippage)
	solUint64 := solLamportBuyAmount.Uint64()
	if err := binary.Write(buf, binary.LittleEndian, solUint64); err != nil {
		return nil, fmt.Errorf("error writing sol lamport buy amount to buffer: %w", err)
	}

	return buf.Bytes(), nil
}
