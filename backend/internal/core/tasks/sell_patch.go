package tasks

import (
	"fmt"
	"personal_bot/internal/core/wallets"

	"github.com/gagliardetto/solana-go"
)

type SellPatch struct {
	Wallet      *wallets.SolanaWallet
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
		st.Wallet = p.Wallet.PrivateKey
		st.WalletName = p.Wallet.WalletName
	}
	if p.Token != nil {
		st.Token = *p.Token
	}
	if p.Amount != nil {
		st.SellPercentage = *p.Amount
	}
	if p.Fee != nil {
		st.Fee = *p.Fee
	}
	if p.Slippage != nil {
		st.Slippage = *p.Slippage
	}
	if p.ComputeUnit != nil {
		st.ComputeUnits = *p.ComputeUnit
	}
	return nil

}
