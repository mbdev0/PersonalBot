package tasks

import (
	"fmt"
	"math/big"
	"pump_fun/pkg/logger"
	"slices"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
)

type Task interface {
	Id() string
	InitDefaults()
}

type BuyTask struct {
	TaskId       string
	Wallet       solana.PrivateKey `validate:"required"`
	TokenAddress solana.PublicKey  `validate:"required"`
	BuyAmount    big.Int           `validate:"required,gtZero"`
	BuyFee       float64           `validate:"required,gt=0"`
	Slippage     float64           `validate:"required,gt=0,lt=1"` // Slippage percentage (0.0 to 1.0)
	ComputeUnits uint32            `validate:"required,min=1"`
	State        State
}

func (bt *BuyTask) InitDefaults() {
	bt.TaskId = uuid.NewString()
	bt.State.SetState(TaskStateCreated)
}

func (bt *BuyTask) TransitionToNextState(nextState TaskState) error {
	validTransitions, ok := StateTransitions[bt.State.TaskState]
	if !ok {
		return fmt.Errorf("invalid current state: %s", bt.State.TaskState.ToString())
	}
	if len(validTransitions) == 0 {
		return fmt.Errorf("no valid transitions from state: %s", bt.State.TaskState.ToString())
	}
	if slices.Contains(validTransitions, nextState) {
		logger.Information(fmt.Sprintf("Transitioning from %s to %s", bt.State.TaskState.ToString(), nextState.ToString()))
		bt.State.SetState(nextState)
		return nil
	}

	return fmt.Errorf("invalid state transition from %s to %s", bt.State.TaskState.ToString(), nextState.ToString())
}

func (bt *BuyTask) Id() string {
	return bt.TaskId
}
