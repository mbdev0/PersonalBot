package tasks

import (
	"math/big"
	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
)

type BuyTask struct {
	taskType     string
	id           string
	Wallet       solana.PrivateKey `validate:"required"`
	Token        solana.PublicKey  `validate:"required"`
	BuyAmount    *big.Int          `validate:"required,gtZero"`
	Fee          float64           `validate:"required,gt=0"`
	Slippage     float64           `validate:"required,gt=0,lt=1"` // Slippage percentage (0.0 to 1.0)
	ComputeUnits uint32            `validate:"required,min=1"`
	state        State
	mu           *sync.RWMutex
}

func NewBuyTask(pk solana.PrivateKey, token solana.PublicKey, common []Option, buyOpts []BuyOption) *BuyTask {
	// ctx, cancel := context.WithCancel(context.Background())
	bt := &BuyTask{
		taskType: "Buy",
		id:       uuid.NewString(),
		Wallet:   pk,
		Token:    token,
		state:    State{TaskState: TaskCreate},
		mu:       &sync.RWMutex{},
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
	bt.mu.RLock()
	defer bt.mu.RUnlock()
	return bt.id
}

func (bt *BuyTask) Type() string { return bt.taskType }                                    //could remove rlock
func (bt *BuyTask) State() State { bt.mu.RLock(); defer bt.mu.RUnlock(); return bt.state } // keep as RLock as many routines might alter

func (bt *BuyTask) SetComputeUnit(cu uint32)     { bt.ComputeUnits = cu } //only main thread sets this
func (bt *BuyTask) SetSlippage(slippage float64) { bt.Slippage = slippage }
func (bt *BuyTask) SetState(newState State)      { bt.mu.Lock(); defer bt.mu.Unlock(); bt.state = newState } //should be locked
