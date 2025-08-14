package tasks

import (
	"context"
	"math/big"
	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
)

type BuyTask struct {
	taskType    string
	id          string
	wallet      solana.PrivateKey `validate:"required"`
	token       solana.PublicKey  `validate:"required"`
	amount      *big.Int          `validate:"required,gtZero"`
	fee         float64           `validate:"required,gt=0"`
	slippage    float64           `validate:"required,gt=0,lt=1"` // Slippage percentage (0.0 to 1.0)
	computeUnit uint32            `validate:"required,min=1"`
	state       State
	ctx         context.Context
	cancel      context.CancelFunc
	mu          sync.RWMutex
}

func NewBuyTask(pk solana.PrivateKey, token solana.PublicKey, common []Option, buyOpts []BuyOption) *BuyTask {
	ctx, cancel := context.WithCancel(context.Background())
	bt := &BuyTask{
		taskType: "Buy",
		id:       uuid.NewString(),
		wallet:   pk,
		token:    token,
		state:    State{TaskState: TaskCreate},
		ctx:      ctx,
		cancel:   cancel,
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

func (bt *BuyTask) Type() string              { bt.mu.RLock(); defer bt.mu.RUnlock(); return bt.taskType }
func (bt *BuyTask) ComputeUnits() uint32      { bt.mu.RLock(); defer bt.mu.RUnlock(); return bt.computeUnit }
func (bt *BuyTask) Slippage() float64         { bt.mu.RLock(); defer bt.mu.RUnlock(); return bt.slippage }
func (bt *BuyTask) Token() solana.PublicKey   { bt.mu.RLock(); defer bt.mu.RUnlock(); return bt.token }
func (bt *BuyTask) Fee() float64              { bt.mu.RLock(); defer bt.mu.RUnlock(); return bt.fee }
func (bt *BuyTask) BuyAmount() *big.Int       { bt.mu.RLock(); defer bt.mu.RUnlock(); return bt.amount }
func (bt *BuyTask) Wallet() solana.PrivateKey { bt.mu.RLock(); defer bt.mu.RUnlock(); return bt.wallet }
func (bt *BuyTask) State() State              { bt.mu.RLock(); defer bt.mu.RUnlock(); return bt.state }

func (bt *BuyTask) SetComputeUnit(cu uint32) { bt.mu.Lock(); defer bt.mu.Unlock(); bt.computeUnit = cu }
func (bt *BuyTask) SetSlippage(slippage float64) {
	bt.mu.Lock()
	defer bt.mu.Unlock()
	bt.slippage = slippage
}
func (bt *BuyTask) SetState(newState State) { bt.mu.Lock(); defer bt.mu.Unlock(); bt.state = newState }

func (bt *BuyTask) ResetCtx() {
	bt.mu.Lock()
	bt.ctx, bt.cancel = context.WithCancel(context.Background())
	bt.mu.Unlock()
}

func (bt *BuyTask) Cancel() {
	if bt.ctx != nil {
		bt.cancel()
	}
}

func (bt *BuyTask) Ctx() context.Context { bt.mu.RLock(); defer bt.mu.RUnlock(); return bt.ctx }
