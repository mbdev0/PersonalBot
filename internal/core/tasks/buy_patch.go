package tasks

import (
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"
)

type BuyPatch struct {
	Wallet      *solana.PrivateKey
	Token       *solana.PublicKey
	Amount      *big.Int
	Fee         *float64
	Slippage    *float64
	ComputeUnit *uint32
}

func (p *BuyPatch) ApplyTo(t Task) error {
	bt, ok := t.(*BuyTask)
	if !ok {
		return fmt.Errorf("patch/task type mismatch: want BuyTask")
	}

	bt.mu.Lock()
	defer bt.mu.Unlock()

	if p.Wallet != nil {
		bt.wallet = *p.Wallet
	}
	if p.Token != nil {
		bt.token = *p.Token
	}
	if p.Amount != nil {
		bt.amount = new(big.Int).Set(p.Amount)
	}
	if p.Fee != nil {
		bt.fee = *p.Fee
	}
	if p.Slippage != nil {
		bt.slippage = *p.Slippage
	}
	if p.ComputeUnit != nil {
		bt.computeUnit = *p.ComputeUnit
	}
	return nil
}
