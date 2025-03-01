package transactions

import (
	"pump_fun/internal/constants"

	"github.com/mr-tron/base58"
)

func GetBondingCurveAddress(tokenAddress string) (string, error) {
	caBytes, _ := base58.Decode(tokenAddress)
	programId, _ := base58.Decode(constants.ProgramID)
	seeds := [][]byte{[]byte("bonding-curve"), caBytes}
	address, _, err := FindProgramAddressSync(seeds, programId)
	return address, err
}
