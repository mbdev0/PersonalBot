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
	Position_id    string
	Wallet         solana.PrivateKey
	Token          solana.PublicKey
	SellPercentage float64
	Fee            float64
	Slippage       float64
	ComputeUnits   uint32
	state          State
	ctx            context.Context
	cancel         context.CancelFunc
	mu             *sync.RWMutex
}

func NewSellTask(pk solana.PrivateKey, token solana.PublicKey, common []Option, sellOpt []SellOption) *SellTask {

	st := &SellTask{
		id:       uuid.NewString(),
		taskType: "Sell",
		state:    State{TaskState: TaskCreate},
		Wallet:   pk,
		Token:    token,
		mu:       &sync.RWMutex{},
	}

	for _, opt := range common {
		opt(st)
	}

	for _, stOpt := range sellOpt {
		stOpt(st)
	}

	return st
}

func (st *SellTask) State() State { st.mu.RLock(); defer st.mu.RUnlock(); return st.state }
func (st *SellTask) Id() string   { return st.id }
func (st *SellTask) Type() string { return st.taskType }

func (st *SellTask) SetComputeUnit(cu uint32)     { st.ComputeUnits = cu }
func (st *SellTask) SetSlippage(slippage float64) { st.Slippage = slippage }
func (st *SellTask) SetState(newState State) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.state = newState
}
