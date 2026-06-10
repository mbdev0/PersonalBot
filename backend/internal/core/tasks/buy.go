package tasks

import (
	"math/big"
	"personal_bot/internal/core/wallets"

	"github.com/gagliardetto/solana-go/rpc"

	"sync"

	"github.com/gagliardetto/solana-go"
)

type BuyTask struct {
	program        string
	taskType       string
	id             int64
	WalletName     string
	Wallet         solana.PrivateKey `validate:"required"`
	WalletId       string
	Token          solana.PublicKey `validate:"required"`
	BuyAmount      *big.Int         `validate:"required,gtZero"`
	Fee            float64          `validate:"required,gt=0"`
	Slippage       float64          `validate:"required,gt=0,lt=1"` // Slippage percentage (0.0 to 1.0)
	ComputeUnits   uint32           `validate:"required,min=1"`
	StrategyId     *int64
	TimeCreated    int64
	state          State
	message        string
	mu             *sync.RWMutex
	rpcGroupId     int64
	rpcNode        *rpc.Client
	rpcNodeString  string
	ws             string
	retries        uint16
	retriesDelayMs uint32
}

func NewBuyTask(wallet wallets.SolanaWallet, token solana.PublicKey, common []Option, buyOpts []BuyOption) *BuyTask {
	bt := &BuyTask{
		taskType:   "Buy",
		WalletId:   wallet.Id,
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

func (bt *BuyTask) Program() string {
	return bt.program
}

func (bt *BuyTask) SetProgram(program string) {
	bt.program = program
}

func (bt *BuyTask) Id() int64 {
	return bt.id
}
func (bt *BuyTask) SetId(id int64) {
	bt.id = id
}

func (bt *BuyTask) GetStrategyId() *int64 {
	return bt.StrategyId
}

func (bt *BuyTask) Type() string                 { return bt.taskType }
func (bt *BuyTask) State() State                 { bt.mu.RLock(); defer bt.mu.RUnlock(); return bt.state } // keep as RLock as many routines might alter
func (bt *BuyTask) SetComputeUnit(cu uint32)     { bt.ComputeUnits = cu }
func (bt *BuyTask) SetSlippage(slippage float64) { bt.Slippage = slippage }
func (bt *BuyTask) SetState(newState State)      { bt.mu.Lock(); defer bt.mu.Unlock(); bt.state = newState }
func (bt *BuyTask) SetTime(t int64)              { bt.TimeCreated = t }
func (bt *BuyTask) SetStrategyId(id int64)       { bt.mu.Lock(); defer bt.mu.Unlock(); bt.StrategyId = &id }
func (bt *BuyTask) SetMessage(message string) {
	bt.message = message
}
func (bt *BuyTask) Message() string {
	return bt.message
}
func (bt *BuyTask) GetWallet() string {
	return bt.Wallet.PublicKey().String()
}

func (bt *BuyTask) GetToken() string {
	return bt.Token.String()
}

func (bt *BuyTask) HttpClient() *rpc.Client {
	return bt.rpcNode
}

func (bt *BuyTask) HttpNode() string {
	return bt.rpcNodeString
}

func (bt *BuyTask) SetHttpNode(rpcNode string) {
	bt.rpcNodeString = rpcNode
	bt.rpcNode = rpc.New(rpcNode)
}

func (bt *BuyTask) WSNode() string {
	return bt.ws
}

func (bt *BuyTask) SetWSNode(ws string) {
	bt.ws = ws
}

func (bt *BuyTask) GetRPCGroupId() int64 {
	return bt.rpcGroupId
}
func (bt *BuyTask) SetRPCGroupId(id int64) {
	bt.rpcGroupId = id
}

func (bt *BuyTask) SetRetries(retries uint16) {
	bt.retries = retries
}

func (bt *BuyTask) SetRetriesDelayMS(retriesDelayMs uint32) {
	bt.retriesDelayMs = retriesDelayMs
}

func (bt *BuyTask) Retries() uint16 { return bt.retries }

func (bt *BuyTask) RetriesDelayMS() uint32 { return bt.retriesDelayMs }

func (bt *BuyTask) RetryFrom() TaskState { return TxInstructionBuild }
