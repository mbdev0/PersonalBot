package tasks

import (
	"fmt"
	"math/big"
	"personal_bot/internal/core/models/wallets"

	"github.com/gagliardetto/solana-go"
)

type BuyPatch struct {
	Wallet      *wallets.SolanaWallet
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
		bt.Wallet = p.Wallet.PrivateKey
		bt.WalletName = p.Wallet.WalletName
	}

	if p.Token != nil {
		bt.Token = *p.Token
	}
	if p.Amount != nil {
		bt.BuyAmount = new(big.Int).Set(p.Amount)
	}
	if p.Fee != nil {
		bt.Fee = *p.Fee
	}
	if p.Slippage != nil {
		bt.Slippage = *p.Slippage
	}
	if p.ComputeUnit != nil {
		bt.ComputeUnits = *p.ComputeUnit
	}
	return nil
}
