package tasks

import (
	"math/big"
	"pump_fun/internal/core/models/wallets"
	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
)

type BuyTask struct {
	taskType     string
	id           string
	WalletName   string
	Wallet       solana.PrivateKey `validate:"required"`
	Token        solana.PublicKey  `validate:"required"`
	BuyAmount    *big.Int          `validate:"required,gtZero"`
	Fee          float64           `validate:"required,gt=0"`
	Slippage     float64           `validate:"required,gt=0,lt=1"` // Slippage percentage (0.0 to 1.0)
	ComputeUnits uint32            `validate:"required,min=1"`
	state        State
	mu           *sync.RWMutex
}

func NewBuyTask(wallet wallets.SolanaWallet, token solana.PublicKey, common []Option, buyOpts []BuyOption) *BuyTask {
	bt := &BuyTask{
		taskType:   "Buy",
		id:         uuid.NewString(),
		WalletName: wallet.WalletName,
		Wallet:     wallet.PrivateKey,
		Token:      token,
		state:      State{TaskState: TaskCreate},
		mu:         &sync.RWMutex{},
	}

	for _, opts := range common {
		opts(bt)
	}

	for _, buyOp := range buyOpts {
		buyOp(bt)
	}

	return bt
}

func (bt *BuyTask) Id() string {
	return bt.id
}

func (bt *BuyTask) Type() string { return bt.taskType }
func (bt *BuyTask) State() State { bt.mu.RLock(); defer bt.mu.RUnlock(); return bt.state } // keep as RLock as many routines might alter

func (bt *BuyTask) SetComputeUnit(cu uint32)     { bt.ComputeUnits = cu }
func (bt *BuyTask) SetSlippage(slippage float64) { bt.Slippage = slippage }
func (bt *BuyTask) SetState(newState State)      { bt.mu.Lock(); defer bt.mu.Unlock(); bt.state = newState }
