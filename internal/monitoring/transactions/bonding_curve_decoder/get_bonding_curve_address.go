package bonding_curve_decoder

import (
	"pump_fun/internal/constants"
	"pump_fun/internal/monitoring/transactions/program_derived_address"

	"github.com/mr-tron/base58"
)

/*

EXAMPLE
address, err := bonding_curve_decoder.GetBondingCurveAddress(coin.CoinData.TokenAddr)
if err != nil {
	fmt.Println("ERROR FINDING PROGRAM ADDRESS\n", err)
} else {
	fmt.Println(address)
}
*/

func GetBondingCurveAddress(tokenAddress string) (string, error) {
	caBytes, _ := base58.Decode(tokenAddress)
	programId, _ := base58.Decode(constants.ProgramID)
	seeds := [][]byte{[]byte("bonding-curve"), caBytes}
	address, _, err := program_derived_address.FindProgramAddressSync(seeds, programId)
	return address, err
}
