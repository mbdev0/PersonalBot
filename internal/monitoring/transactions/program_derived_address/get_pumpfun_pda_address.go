package program_derived_address

import (
	"pump_fun/internal/constants"

	"github.com/mr-tron/base58"
)

func GetBondingCurveAddress(tokenAddress string) (bondingCurveAddress string, err error) {
	caBytes, _ := base58.Decode(tokenAddress)
	programId, _ := base58.Decode(constants.Program)
	seeds := [][]byte{[]byte("bonding-curve"), caBytes}
	address, _, err := FindProgramAddressSync(seeds, programId)
	return address, err
}

func GetAssociatedBondingCurveAddress(bondingCurveAddr string, tokenAddress string) (associatedBondingCurveAddress string, err error) {
	caBytes, _ := base58.Decode(tokenAddress)
	tokenProgram, _ := base58.Decode(constants.TokenProgram)
	bondingCurve, _ := base58.Decode(bondingCurveAddr)
	associatedTokenProgram, _ := base58.Decode(constants.AssociatedTokenProgram)

	seeds := [][]byte{bondingCurve, tokenProgram, caBytes}
	address, _, err := FindProgramAddressSync(seeds, associatedTokenProgram)
	return address, err
}

func GetCreatorVaultAddress(devAddress string) (creatorVaultAddress string, err error) {
	caBytes, _ := base58.Decode(devAddress)
	programId, _ := base58.Decode(constants.Program)
	seeds := [][]byte{[]byte("creator-vault"), caBytes}
	address, _, err := FindProgramAddressSync(seeds, programId)
	return address, err
}
