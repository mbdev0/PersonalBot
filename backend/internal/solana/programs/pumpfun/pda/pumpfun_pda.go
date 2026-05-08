package pda

import (
	"personal_bot/internal/solana/utils"

	"github.com/mr-tron/base58"
)

func GetBondingCurveAddress(tokenAddress string, program string) (bondingCurveAddress string, err error) {
	caBytes, _ := base58.Decode(tokenAddress)
	programId, _ := base58.Decode(program)
	seeds := [][]byte{[]byte("bonding-curve"), caBytes}
	address, _, err := utils.FindProgramAddressSync(seeds, programId)

	return address, err
}
