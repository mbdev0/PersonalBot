package strategies

import (
	"fmt"
	"math/big"
	rpcgroups "personal_bot/internal/core/rpc_groups"
	"personal_bot/internal/core/wallets"
)

type AfkPatch struct {
	Filters        *[]StrategyFilter
	BuyAmount      *big.Int
	BuyFee         *float64
	Slippage       *float64
	ComputeUnits   *float64
	Wallet         *wallets.SolanaWallet
	SellStrategies *[]StrategyConfig
	SellFee        *float64
	RPCGroup       *rpcgroups.RPCGroup
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

	if ap.RPCGroup != nil {
		afk.RPCGroup = *ap.RPCGroup
	}

	afk.SellFee = ap.SellFee

	return nil

}
