package pda

import (
	"personal_bot/internal/core/constants"
	"personal_bot/internal/solana/utils"

	"github.com/mr-tron/base58"
)

func GetBondingCurveAddress(tokenAddress string) (bondingCurveAddress string, err error) {
	caBytes, _ := base58.Decode(tokenAddress)
	programId, _ := base58.Decode(constants.PumpFunProgram)
	seeds := [][]byte{[]byte("bonding-curve"), caBytes}
	address, _, err := utils.FindProgramAddressSync(seeds, programId)

	return address, err
}

func GetBondingCurveV2Address(tokenAddress string) (bondingCurveV2Address string, err error) {
	caBytes, _ := base58.Decode(tokenAddress)
	programId, _ := base58.Decode(constants.PumpFunProgram)
	seeds := [][]byte{[]byte("bonding-curve-v2"), caBytes}
	address, _, err := utils.FindProgramAddressSync(seeds, programId)

	return address, err
}

func GetAssociatedBondingCurveAddress(bondingCurveAddr string, tokenAddress string, isNewTokenAddress bool) (associatedBondingCurveAddress string, err error) {
	caBytes, _ := base58.Decode(tokenAddress)

	var tokenProgramAddress string
	if isNewTokenAddress {
		tokenProgramAddress = constants.Token2022Program
	} else {
		tokenProgramAddress = constants.TokenProgram
	}

	tokenProgram, _ := base58.Decode(tokenProgramAddress)
	bondingCurve, _ := base58.Decode(bondingCurveAddr)
	associatedTokenProgram, _ := base58.Decode(constants.AssociatedTokenProgram)

	seeds := [][]byte{bondingCurve, tokenProgram, caBytes}
	address, _, err := utils.FindProgramAddressSync(seeds, associatedTokenProgram)
	return address, err
}

func GetCreatorVaultAddress(devAddress string) (creatorVaultAddress string, err error) {
	caBytes, _ := base58.Decode(devAddress)
	programId, _ := base58.Decode(constants.PumpFunProgram)
	seeds := [][]byte{[]byte("creator-vault"), caBytes}
	address, _, err := utils.FindProgramAddressSync(seeds, programId)
	return address, err
}

func GetUserVolumeAccumulatorAddress(walletAddress string) (userVolumeAccumulatorAddress string, err error) {
	walletBytes, _ := base58.Decode(walletAddress)
	programId, _ := base58.Decode(constants.PumpFunProgram)
	seeds := [][]byte{[]byte("user_volume_accumulator"), walletBytes}
	address, _, err := utils.FindProgramAddressSync(seeds, programId)

	if err != nil {
		return "", err
	}

	return address, nil
}
