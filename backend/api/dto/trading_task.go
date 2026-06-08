package dto

import (
	programnames "personal_bot/app/program_names"

	"github.com/go-playground/validator/v10"
)

type TradingType string

const (
	AFK  TradingType = "AFK"
	BUY  TradingType = "BUY"
	SELL TradingType = "SELL"
	SPAM TradingType = "SPAM"
)

type TradingTask struct {
	Program      programnames.Program `json:"program" validate:"required,oneof=PumpfunNative PumpfunAMM"`
	Type         TradingType          `json:"trading_type" validate:"required,oneof=AFK BUY SELL SPAM"`
	ComputeUnits float64              `json:"compute_units" validate:"required,gt=0"`
	WalletName   string               `json:"wallet_name" validate:"required"`
	Slippage     float64              `json:"slippage" validate:"required,gte=0,lte=1"`
	RPCGroupId   int64                `json:"rpc_group_id" validate:"required"`

	Filters          *Filters           `json:"filters,omitempty" validate:"required_if=Type AFK"`
	SellStrategies   *[]SellStrategyDTO `json:"sell_strategies,omitempty" validate:"required_if=Type AFK,required_if=Type BUY"`
	BuyAmount        *float64           `json:"buy_amount,omitempty" validate:"required_if=Type AFK,required_if=Type BUY"`
	BuyFee           *float64           `json:"buy_fee,omitempty" validate:"required_if=Type AFK,required_if=Type BUY"`
	SellAmount       *float64           `json:"sell_amount,omitempty" validate:"required_if=Type Sell"`
	SellFee          *float64           `json:"sell_fee,omitempty"`
	TokenAddress     *string            `json:"token_address,omitempty" validate:"required_if=Type SELL,required_if=Type BUY"`
	Retries          *uint16            `json:"retries,omitempty" validate:"omitnil,gte=0"`
	RetriesDelayMS   *uint32            `json:"retry_delay_ms,omitempty" validate:"omitnil,gte=0"`
	NumberOfSubTasks *uint32            `json:"no_of_tasks,omitempty" validate:"omitnil,gte=0,required_if=Type SPAM"`
	StartTime        *uint64            `json:"start_time,omitempty" validate:"omitnil,gte=0,required_if=Type SPAM"`
}

func (t *TradingTask) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}

type TradingTaskResponse struct {
	Program       string      `json:"program"`
	Type          TradingType `json:"trading_type"`
	Id            int64       `json:"id"`
	ComputeUnits  float64     `json:"compute_units"`
	Slippage      float64     `json:"slippage"`
	WalletName    string      `json:"wallet_name"`
	WalletAddress string      `json:"wallet_address"`
	State         string      `json:"state"`
	Message       string      `json:"message"`
	RPCGroup      string      `json:"rpc_group"`
	RPCGroupID    int64       `json:"rpc_group_id"`

	SellStrategies   *[]SellStrategyDTO `json:"sell_strategies,omitempty"`
	BuyAmount        *float64           `json:"buy_amount,omitempty"`
	Filters          *Filters           `json:"filters,omitempty"`
	SellFee          *float64           `json:"sell_fee,omitempty"`
	BuyFee           *float64           `json:"buy_fee,omitempty"`
	SellAmount       *float64           `json:"sell_amount,omitempty"`
	SellTaskId       *int64             `json:"sell_task_id,omitempty"`
	BuyTaskId        *int64             `json:"buy_task_id,omitempty"`
	PositionId       *int64             `json:"position_id,omitempty"`
	TokenAddress     *string            `json:"token_address,omitempty"`
	TimeCreated      int64              `json:"time_created"`
	Retries          *uint16            `json:"retries,omitempty"`
	RetriesDelayMS   *uint32            `json:"retry_delay_ms,omitempty"`
	NumberOfSubTasks *uint32            `json:"no_of_tasks,omitempty"`
	StartTime        *uint64            `json:"start_time,omitempty"`
}

type TradingTaskPatch struct {
	Program      *programnames.Program `json:"program"`
	Type         *TradingType          `json:"trading_type" validate:"omitempty,oneof=AFK BUY SELL"`
	ComputeUnits *float64              `json:"compute_units" validate:"omitempty,gt=0"`
	WalletName   *string               `json:"wallet_name" validate:"omitempty"`
	Slippage     *float64              `json:"slippage" validate:"omitempty,gte=0,lte=1"`
	RPCGroupId   *int64                `json:"rpc_group_id" validate:"omitempty"`

	Filters          *Filters           `json:"filters,omitempty"`
	SellStrategies   *[]SellStrategyDTO `json:"sell_strategies,omitempty"`
	BuyAmount        *float64           `json:"buy_amount,omitempty" validate:"omitempty,gt=0"`
	BuyFee           *float64           `json:"buy_fee,omitempty" validate:"omitempty,gte=0"`
	SellAmount       *float64           `json:"sell_amount,omitempty" validate:"omitempty,gt=0"`
	SellFee          *float64           `json:"sell_fee,omitempty" validate:"omitempty,gte=0"`
	TokenAddress     *string            `json:"token_address,omitempty"`
	Retries          *uint16            `json:"retries,omitempty"`
	RetriesDelayMS   *uint32            `json:"retry_delay_ms,omitempty"`
	NumberOfSubTasks *uint32            `json:"no_of_tasks"`
	StartTime        *uint64            `json:"start_time"`
}

func (t *TradingTaskPatch) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}
