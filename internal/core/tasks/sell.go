package tasks

import (
	"context"
	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
)

type SellTask struct {
	taskType       string
	id             string
	wallet         solana.PrivateKey
	token          solana.PublicKey
	sellPercentage float64
	fee            float64
	slippage       float64
	computeUnits   uint32
	state          State
	ctx            context.Context
	cancel         context.CancelFunc
	mu             sync.RWMutex
}

func NewSellTask(pk solana.PrivateKey, token solana.PublicKey, common []Option, sellOpt []SellOption) *SellTask {
	ctx, cancel := context.WithCancel(context.Background())

	st := &SellTask{
		id:       uuid.NewString(),
		taskType: "Sell",
		state:    State{TaskState: TaskCreate},
		wallet:   pk,
		token:    token,
		ctx:      ctx,
		cancel:   cancel,
	}

	for _, opt := range common {
		opt(st)
	}

	for _, stOpt := range sellOpt {
		stOpt(st)
	}

	return st
}

func (st *SellTask) State() State            { st.mu.RLock(); defer st.mu.RUnlock(); return st.state }
func (st *SellTask) Id() string              { st.mu.RLock(); defer st.mu.RUnlock(); return st.id }
func (st *SellTask) Type() string            { st.mu.RLock(); defer st.mu.RUnlock(); return st.taskType }
func (st *SellTask) Slippage() float64       { st.mu.RLock(); defer st.mu.RUnlock(); return st.slippage }
func (st *SellTask) Token() solana.PublicKey { st.mu.RLock(); defer st.mu.RUnlock(); return st.token }
func (st *SellTask) Fee() float64            { st.mu.RLock(); defer st.mu.RUnlock(); return st.fee }
func (st *SellTask) ComputeUnits() uint32 {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.computeUnits
}

func (st *SellTask) SellPercentage() float64 {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.sellPercentage
}

func (st *SellTask) Wallet() solana.PrivateKey {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.wallet
}

func (st *SellTask) SetComputeUnit(cu uint32) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.computeUnits = cu
}
func (st *SellTask) SetSlippage(slippage float64) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.slippage = slippage
}
func (st *SellTask) SetState(newState State) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.state = newState
}

func (st *SellTask) ResetCtx() {
	ctx, cancel := context.WithCancel(context.Background())
	st.ctx = ctx
	st.cancel = cancel
}

func (st *SellTask) Cancel() {
	if st.ctx != nil {
		st.cancel()

	}
}

func (st *SellTask) Ctx() context.Context {
	return st.ctx
}
