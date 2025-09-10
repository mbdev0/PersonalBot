package coordination

import (
	"fmt"
	"math/big"
	"pump_fun/infrastructure/solana_price"
	positionmodels "pump_fun/internal/core/position"
	"pump_fun/internal/services/position"
	subscriptionhub "pump_fun/internal/services/subscription_hub"
	taskservice "pump_fun/internal/services/task_service"
	"pump_fun/pkg/logger"
)

/*
this will be in charge of just streaming the positions to ws
*/

type PositionReporter struct {
	positionService *position.Service
	taskService     *taskservice.TaskService
	publisher       subscriptionhub.Publisher
}

// this will stream the positions
func NewPositionReporter(posService *position.Service, taskService *taskservice.TaskService, publisher subscriptionhub.Publisher) *PositionReporter {
	positionreporter := PositionReporter{
		positionService: posService,
		taskService:     taskService,
		publisher:       publisher,
	}

	return &positionreporter
}

func (pr *PositionReporter) StreamPosition(marketCapChan chan float64, taskId string) {
	task, err := pr.taskService.GetTaskWith(taskId)
	if err != nil {
		logger.Error(err)
	}

	pos, err := pr.positionService.GetById(taskId)
	if err != nil {
		logger.Error(err)
		// return nil, err
	}

	for mcap := range marketCapChan {
		positionMessage := pr.generatePositionMessage(pos, mcap)
		fmt.Println(positionMessage)
		//publish the profit, mcap, position
		pr.publisher.PublishPositionUpdate(task, positionMessage)
	}
}

func (pr *PositionReporter) generatePositionMessage(pos *positionmodels.Position, marketCap float64) positionmodels.PositionMessage {
	//we need unrealized profit -> remaning tokens (in sol)
	totalPnl, unrealizedPnl := pr.getProfitValues(pos, marketCap)

	posMessage := positionmodels.PositionMessage{
		BuyTaskId:        pos.PositionId,
		UnrealizedProfit: unrealizedPnl.Text('f', 9),
		RealizedProfit:   pos.FinalizedProfit.Text('f', 9),
		TotalPnL:         totalPnl.Text('f', 9),
		MarketCap:        marketCap,
		RemainingTokens:  pos.TokenRemaining.String(),
	}

	return posMessage
}

func (pr *PositionReporter) getProfitValues(pos *positionmodels.Position, marketCap float64) (totalPnl *big.Float, unrealized *big.Float) {
	solPrice, err := solana_price.GetSolPrice()
	if err != nil {
		logger.Error(err)
	}

	tokenValue := pr.calculateTokenValueInSol(marketCap, pos.InitialTokenAmount, *solPrice)
	unrealizedPnl := new(big.Float).Sub(tokenValue, pos.RemaningCostBasis)

	totalPnL := new(big.Float).Add(pos.FinalizedProfit, unrealizedPnl)

	return totalPnL, unrealizedPnl
}

func (pr *PositionReporter) calculateTokenValueInSol(marketCapUSD float64, tokensRemaining *big.Float, solPrice float64) (totalValueSOL *big.Float) {
	marketCapSol := new(big.Float).Quo(big.NewFloat(marketCapUSD), big.NewFloat(solPrice))
	totalSupply := new(big.Float).SetInt64(1000000000)

	pricePerTokenSOL := new(big.Float).Quo(marketCapSol, totalSupply)

	totalValueSOL = new(big.Float).Mul(pricePerTokenSOL, tokensRemaining)

	return totalValueSOL
}
