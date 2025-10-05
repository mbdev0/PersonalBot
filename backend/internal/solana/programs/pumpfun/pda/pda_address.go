package pda

import (
	"pump_fun/internal/core/constants"
	"pump_fun/internal/solana/utils"

	"github.com/mr-tron/base58"
)

func GetBondingCurveAddress(tokenAddress string) (bondingCurveAddress string, err error) {
	caBytes, _ := base58.Decode(tokenAddress)
	programId, _ := base58.Decode(constants.Program)
	seeds := [][]byte{[]byte("bonding-curve"), caBytes}
	address, _, err := utils.FindProgramAddressSync(seeds, programId)
	return address, err
}

func GetAssociatedBondingCurveAddress(bondingCurveAddr string, tokenAddress string) (associatedBondingCurveAddress string, err error) {
	caBytes, _ := base58.Decode(tokenAddress)
	tokenProgram, _ := base58.Decode(constants.TokenProgram)
	bondingCurve, _ := base58.Decode(bondingCurveAddr)
	associatedTokenProgram, _ := base58.Decode(constants.AssociatedTokenProgram)

	seeds := [][]byte{bondingCurve, tokenProgram, caBytes}
	address, _, err := utils.FindProgramAddressSync(seeds, associatedTokenProgram)
	return address, err
}

func GetCreatorVaultAddress(devAddress string) (creatorVaultAddress string, err error) {
	caBytes, _ := base58.Decode(devAddress)
	programId, _ := base58.Decode(constants.Program)
	seeds := [][]byte{[]byte("creator-vault"), caBytes}
	address, _, err := utils.FindProgramAddressSync(seeds, programId)
	return address, err
}

func GetUserVolumeAccumulatorAddress(walletAddress string) (userVolumeAccumulatorAddress string, err error) {
	walletBytes, _ := base58.Decode(walletAddress)
	programId, _ := base58.Decode(constants.Program)
	seeds := [][]byte{[]byte("user_volume_accumulator"), walletBytes}
	address, _, err := utils.FindProgramAddressSync(seeds, programId)

	if err != nil {
		return "", err
	}

	return address, nil
}
