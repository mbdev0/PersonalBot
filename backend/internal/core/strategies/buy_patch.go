package strategies

import (
	"fmt"
	"math/big"
	"personal_bot/internal/core/models/wallets"

	"github.com/gagliardetto/solana-go"
)

type BuyPatch struct {
	BuyAmount      *big.Int
	BuyFee         *float64
	Slippage       *float64
	ComputeUnits   *float64
	Wallet         *wallets.SolanaWallet
	Token          *solana.PublicKey
	SellStrategies *[]StrategyConfig
	SellFee        *float64
}

func (bp *BuyPatch) ApplyTo(task Task) error {
	buy, ok := task.(*Buy)
	if !ok {
		return fmt.Errorf("patch/task type mismatch: want buy")
	}

	if bp.BuyAmount != nil {
		buy.BuyAmount = bp.BuyAmount
	}

	if bp.BuyFee != nil {
		buy.BuyFee = *bp.BuyFee
	}

	if bp.ComputeUnits != nil {
		buy.ComputeUnits = *bp.ComputeUnits
	}

	if bp.Slippage != nil {
		buy.Slippage = *bp.Slippage
	}

	if bp.Wallet != nil {
		buy.Wallet = *bp.Wallet
	}

	if bp.SellStrategies != nil {
		buy.SellStrategies = *bp.SellStrategies
	}

	if bp.Token != nil {
		buy.Token = *bp.Token
	}

	buy.SellFee = bp.SellFee

	return nil

}
