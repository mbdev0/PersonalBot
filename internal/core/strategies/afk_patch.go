package strategies

import (
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"
)

type AfkPatch struct {
	Filters        *[]StrategyFilter
	BuyAmount      *big.Int
	BuyFee         *float64
	Slippage       *float64
	ComputeUnits   *float64
	Wallet         *solana.PrivateKey
	SellStrategies *[]StrategyConfig
}

func (ap *AfkPatch) ApplyTo(task Task) error {
	afk, ok := task.(*Afk)
	if !ok {
		return fmt.Errorf("patch/task type mismatch: want Afk")
	}

	if ap.BuyAmount != nil {
		afk.BuyAmount = ap.BuyAmount
	}

	if ap.BuyFee != nil {
		afk.BuyFee = *ap.BuyFee
	}

	if ap.ComputeUnits != nil {
		afk.ComputeUnits = *ap.ComputeUnits
	}

	if ap.Slippage != nil {
		afk.Slippage = *ap.Slippage
	}

	if ap.Wallet != nil {
		afk.Wallet = *ap.Wallet
	}

	if ap.Filters != nil {
		afk.Filters = *ap.Filters
	}

	if ap.SellStrategies != nil {
		afk.SellStrategies = *ap.SellStrategies
	}

	return nil

}
