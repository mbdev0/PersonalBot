package tasks

import (
	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
)

type SellTask struct {
	TaskId           string
	Wallet           solana.PrivateKey
	TokenAddress     solana.PublicKey
	PercentageToSell float64
	SellFee          float64
	Slippage         float64
	ComputeUnits     uint32
	State            State
}

func (sellTask *SellTask) Id() string {
	return sellTask.TaskId
}

func (sellTask *SellTask) InitDefaults() {
	sellTask.TaskId = uuid.NewString()
	sellTask.State.TaskState = TaskStateCreated
}
