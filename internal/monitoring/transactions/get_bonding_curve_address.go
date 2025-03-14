package transactions

import (
	"pump_fun/internal/constants"

	"github.com/mr-tron/base58"
)

/*

EXAMPLE
address, err := GetBondingCurveAddress(coin.CoinData.TokenAddr)
if err != nil {
	fmt.Println("ERROR FINDING PROGRAM ADDRESS\n", err)
} else {
	fmt.Println("BONDING CURVE CHECK: ", coin.CoinData.TokenAddr, coin.CoinData.BondingCurveAddr, coin.CoinData.BondingCurveAddr == address)
}
*/

func GetBondingCurveAddress(tokenAddress string) (string, error) {
	caBytes, _ := base58.Decode(tokenAddress)
	programId, _ := base58.Decode(constants.ProgramID)
	seeds := [][]byte{[]byte("bonding-curve"), caBytes}
	address, _, err := FindProgramAddressSync(seeds, programId)
	return address, err
}
