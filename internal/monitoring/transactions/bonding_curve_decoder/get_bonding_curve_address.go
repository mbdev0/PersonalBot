package bonding_curve_decoder

import (
	"pump_fun/internal/constants"
	"pump_fun/internal/monitoring/transactions/program_derived_address"

	"github.com/mr-tron/base58"
)

func GetBondingCurveAddress(tokenAddress string) (bondingCurveAddress string, err error) {
	caBytes, _ := base58.Decode(tokenAddress)
	programId, _ := base58.Decode(constants.Program)
	seeds := [][]byte{[]byte("bonding-curve"), caBytes}
	address, _, err := program_derived_address.FindProgramAddressSync(seeds, programId)
	return address, err
}

func GetAssociatedBondingCurveAddress(bondingCurveAddr string, tokenAddress string) (associatedBondingCurveAddress string, err error) {
	caBytes, _ := base58.Decode(tokenAddress)
	tokenProgram, _ := base58.Decode(constants.TokenProgram)
	bondingCurve, _ := base58.Decode(bondingCurveAddr)
	associatedTokenProgram, _ := base58.Decode(constants.AssociatedTokenProgram)

	seeds := [][]byte{bondingCurve, tokenProgram, caBytes}
	address, _, err := program_derived_address.FindProgramAddressSync(seeds, associatedTokenProgram)
	return address, err
}
