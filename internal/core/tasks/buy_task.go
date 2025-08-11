package tasks

import (
	"context"
	"math/big"
	"pump_fun/api/dto"
	"pump_fun/internal/core/models"
	"pump_fun/internal/utils"
	"pump_fun/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
)

type Task interface {
	Id() string
	InitDefaults()
	InitCancelToken()
	UpdateTask(newTask dto.RequestTask) (err error)
	SetState(State)
	GetTaskType() string
	GetState() State
	Cancel()
	Context() context.Context
}

type BuyTask struct {
	TaskType     string
	TaskId       string
	Wallet       solana.PrivateKey `validate:"required"`
	TokenAddress solana.PublicKey  `validate:"required"`
	BuyAmount    big.Int           `validate:"required,gtZero"`
	BuyFee       float64           `validate:"required,gt=0"`
	Slippage     float64           `validate:"required,gt=0,lt=1"` // Slippage percentage (0.0 to 1.0)
	ComputeUnits uint32            `validate:"required,min=1"`
	State        State
	CancelToken  models.CancelToken
}

func (bt *BuyTask) InitDefaults() {
	bt.TaskId = uuid.NewString()
	bt.TaskType = "Buy"
	bt.State.TaskState = TaskCreate
}

func (bt *BuyTask) Id() string {
	return bt.TaskId
}

func (bt *BuyTask) UpdateTask(newTask dto.RequestTask) (err error) {
	bt.Wallet, err = solana.PrivateKeyFromBase58(newTask.WalletAddressPrivateKey)
	if err != nil {
		return err
	}

	bt.TokenAddress, err = solana.PublicKeyFromBase58(newTask.TokenAddress)
	if err != nil {
		return err
	}

	bt.BuyAmount = utils.ConvertSolToLamport(*newTask.BuyAmount)
	bt.BuyFee = *newTask.BuyFee
	bt.Slippage = newTask.Slippage
	bt.ComputeUnits = newTask.ComputeUnits

	return nil
}

func (bt *BuyTask) SetState(newState State) {
	bt.State = newState
}

func (bt *BuyTask) GetTaskType() string {
	return bt.TaskType
}

func (bt *BuyTask) GetState() State {
	return bt.State
}

func (bt *BuyTask) InitCancelToken() {
	bt.CancelToken.CancellationContext, bt.CancelToken.CancellationFunc = context.WithCancel(context.Background())
}

func (bt *BuyTask) Cancel() {
	logger.Information("cancel called")
	if bt.CancelToken.CancellationContext != nil {
		bt.CancelToken.CancellationFunc()
	}
}

func (bt *BuyTask) Context() context.Context {
	return bt.CancelToken.CancellationContext
}
