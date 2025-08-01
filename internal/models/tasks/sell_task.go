package tasks

import (
	"pump_fun/api/models"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
)

type SellTask struct {
	TaskType         string
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
	sellTask.TaskType = "Sell"
	sellTask.State.TaskState = TaskStateCreated
}

func (sellTask *SellTask) UpdateTask(newTask models.RequestTask) (err error) {
	sellTask.Wallet, err = solana.PrivateKeyFromBase58(newTask.WalletAddressPrivateKey)
	if err != nil {
		return err
	}

	sellTask.TokenAddress, err = solana.PublicKeyFromBase58(newTask.TokenAddress)
	if err != nil {
		return err
	}

	sellTask.PercentageToSell = *newTask.SellAmount
	sellTask.SellFee = *newTask.SellFee
	sellTask.Slippage = newTask.Slippage
	sellTask.ComputeUnits = newTask.ComputeUnits
	return nil
}

func (sellTask *SellTask) SetState(newState TaskState) {
	sellTask.State.TaskState = newState
}

func (sellTask *SellTask) GetTaskType() string {
	return sellTask.TaskType
}
