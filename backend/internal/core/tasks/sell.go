package tasks

import (
	"personal_bot/internal/core/models/wallets"
	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type SellTask struct {
	taskType       string
	id             int64
	Position_id    *int64
	WalletName     string
	Wallet         solana.PrivateKey
	WalletID       string
	Token          solana.PublicKey
	SellPercentage float64
	Fee            float64
	Slippage       float64
	ComputeUnits   uint32
	TimeCreated    int64
	state          State
	message        string
	StrategyId     *int64
	mu             *sync.RWMutex
	rpcNode        *rpc.Client
	rpcNodeString  string
	ws             string
}

func NewSellTask(wallet wallets.SolanaWallet, token solana.PublicKey, common []Option, sellOpt []SellOption) *SellTask {

	st := &SellTask{
		taskType:   "Sell",
		state:      State{TaskState: TaskCreate},
		WalletName: wallet.WalletName,
		Wallet:     wallet.PrivateKey,
		WalletID:   wallet.Id,
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

func (st *SellTask) GetStrategyId() *int64 {
	return st.StrategyId
}

func (st *SellTask) State() State   { st.mu.RLock(); defer st.mu.RUnlock(); return st.state }
func (st *SellTask) Id() int64      { return st.id }
func (st *SellTask) SetId(id int64) { st.id = id }
func (st *SellTask) Type() string   { return st.taskType }

func (st *SellTask) SetComputeUnit(cu uint32)     { st.ComputeUnits = cu }
func (st *SellTask) SetSlippage(slippage float64) { st.Slippage = slippage }
func (st *SellTask) SetState(newState State) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.state = newState
}
func (st *SellTask) SetTime(t int64)        { st.TimeCreated = t }
func (st *SellTask) SetStrategyId(id int64) { st.mu.Lock(); defer st.mu.Unlock(); st.StrategyId = &id }

func (st *SellTask) SetMessage(message string) {
	st.message = message
}

func (st *SellTask) Message() string {
	return st.message
}

func (st *SellTask) GetWallet() string {
	return st.Wallet.PublicKey().String()
}

func (st *SellTask) GetToken() string {
	return st.Token.String()
}

func (st *SellTask) HttpClient() *rpc.Client {
	return st.rpcNode
}

func (st *SellTask) SetHttpNode(rpcNode string) {
	st.rpcNodeString = rpcNode
	st.rpcNode = rpc.New(rpcNode)
}
func (st *SellTask) HttpNode() string {
	return st.rpcNodeString
}

func (st *SellTask) WSNode() string {
	return st.ws
}

func (st *SellTask) SetWSNode(ws string) {
	st.ws = ws
}
