package strategies

import (
	"fmt"
	"personal_bot/internal/core/models/wallets"

	"github.com/gagliardetto/solana-go"
)

type SellPatch struct {
	SellAmount   *float64
	SellFee      *float64
	Slippage     *float64
	ComputeUnits *float64
	Token        *solana.PublicKey
	Wallet       *wallets.SolanaWallet
}

func (sp *SellPatch) ApplyTo(task Task) error {
	sell, ok := task.(*Sell)
	if !ok {
		return fmt.Errorf("patch/task type mismatch: want sell")
	}

	if sp.SellAmount != nil {
		sell.SellAmount = *sp.SellAmount
	}

	if sp.SellFee != nil {
		sell.SellFee = *sp.SellFee
	}

	if sp.ComputeUnits != nil {
		sell.ComputeUnits = *sp.ComputeUnits
	}

	if sp.Slippage != nil {
		sell.Slippage = *sp.Slippage
	}

	if sp.Wallet != nil {
		sell.Wallet = *sp.Wallet
	}

	if sp.Token != nil {
		sell.Token = *sp.Token
	}

	if sp.SellFee != nil {
		sell.SellFee = *sp.SellFee
	}

	return nil

}
