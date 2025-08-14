package tasks

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
)

type SellPatch struct {
	Wallet      *solana.PrivateKey
	Token       *solana.PublicKey
	Amount      *float64
	Fee         *float64
	Slippage    *float64
	ComputeUnit *uint32
}

func (p *SellPatch) ApplyTo(t Task) error {
	st, ok := t.(*SellTask)
	if !ok {
		return fmt.Errorf("patch/task type mismatch: want BuyTask")
	}

	st.mu.Lock()
	defer st.mu.Unlock()

	if p.Wallet != nil {
		st.wallet = *p.Wallet
	}
	if p.Token != nil {
		st.token = *p.Token
	}
	if p.Amount != nil {
		st.sellPercentage = *p.Amount
	}
	if p.Fee != nil {
		st.fee = *p.Fee
	}
	if p.Slippage != nil {
		st.slippage = *p.Slippage
	}
	if p.ComputeUnit != nil {
		st.computeUnits = *p.ComputeUnit
	}
	return nil

}
