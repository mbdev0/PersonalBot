package tasks

import (
	"pump_fun/app/iterable"
	"pump_fun/internal/core/models/wallets"
	"sync"

	"github.com/gagliardetto/solana-go"
)

type SellTask struct {
	taskType       string
	id             int64
	Position_id    *int64
	WalletName     string
	Wallet         solana.PrivateKey
	Token          solana.PublicKey
	SellPercentage float64
	Fee            float64
	Slippage       float64
	ComputeUnits   uint32
	state          State
	mu             *sync.RWMutex
}

func NewSellTask(wallet wallets.SolanaWallet, token solana.PublicKey, common []Option, sellOpt []SellOption) *SellTask {

	st := &SellTask{
		id:         iterable.Itr.ID(),
		taskType:   "Sell",
		state:      State{TaskState: TaskCreate},
		WalletName: wallet.WalletName,
		Wallet:     wallet.PrivateKey,
		Token:      token,
		mu:         &sync.RWMutex{},
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
func (st *SellTask) Id() int64    { return st.id }
func (st *SellTask) Type() string { return st.taskType }

func (st *SellTask) SetComputeUnit(cu uint32)     { st.ComputeUnits = cu }
func (st *SellTask) SetSlippage(slippage float64) { st.Slippage = slippage }
func (st *SellTask) SetState(newState State) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.state = newState
}
