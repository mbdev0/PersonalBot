package position

import (
	"fmt"
	"math/big"
	"pump_fun/infrastructure/solana_price"
	"pump_fun/internal/core/position"
	subscriptionhub "pump_fun/internal/services/subscription_hub"
	taskservice "pump_fun/internal/services/task_service"
	"pump_fun/pkg/logger"
	"sync"

	"github.com/gagliardetto/solana-go"
)

type Service struct {
	positions   map[string]*position.Position
	publisher   subscriptionhub.Publisher
	taskService *taskservice.TaskService
	mu          *sync.Mutex
}

func NewPositionService(subhub subscriptionhub.Publisher, taskService *taskservice.TaskService) Service {
	return Service{
		positions:   map[string]*position.Position{},
		publisher:   subhub, // TODO -> might need to change this?
		taskService: taskService,
		mu:          &sync.Mutex{},
	}
}

func (s *Service) ReportBuy(buytaskid string, tokenaddress solana.PublicKey, walletAddress solana.PublicKey, tokenAmount *big.Int, solSpent *big.Float) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newPosition := position.Position{
		PositionId:         buytaskid,
		TokenAddress:       tokenaddress,
		WalletAddress:      walletAddress,
		InitialTokenAmount: tokenAmount,
		TokenRemaining:     tokenAmount,
		RemaningCostBasis:  solSpent,
		FinalizedProfit:    big.NewFloat(0),
	}

	s.positions[newPosition.PositionId] = &newPosition

	//can send a webhook
	//publish a position has been created for buytaskid
}

func (s *Service) ReportSell(buyTaskId string, tokensSold *big.Int, solRecieved *big.Float) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pos, ok := s.positions[buyTaskId]
	if !ok {
		return fmt.Errorf("position not found with id: %s", buyTaskId)
	}

	if pos.TokenRemaining.Cmp(tokensSold) == -1 {
		return fmt.Errorf("tokens sold is more than the token remaining... something went wrong")
	}

	sellRatio := new(big.Float).Quo(
		new(big.Float).SetInt(tokensSold),
		new(big.Float).SetInt(pos.TokenRemaining),
	)

	costBasisofSoldTokens := new(big.Float).Mul(sellRatio, pos.RemaningCostBasis)
	realizedBasisFromSale := new(big.Float).Sub(solRecieved, costBasisofSoldTokens)

	pos.RemaningCostBasis.Sub(pos.RemaningCostBasis, costBasisofSoldTokens)
	pos.FinalizedProfit.Add(pos.FinalizedProfit, realizedBasisFromSale)
	pos.TokenRemaining.Sub(pos.TokenRemaining, tokensSold)

	//can send a webhook here too
	//publish the position for buytaskid has been updated

	return nil
}

func (s *Service) GetById(id string) (*position.Position, error) {
	pos, ok := s.positions[id]
	if !ok {
		return nil, fmt.Errorf("position not found with id %s", id)
	}

	return pos, nil
}

func (s *Service) GetAll() []position.Position {
	allPos := make([]position.Position, len(s.positions))

	var index int
	for _, pos := range s.positions {
		allPos[index] = *pos
		index += 1
	}

	return allPos
}

func (s *Service) StreamPosition(marketCapChan chan float64, taskId string) {
	task, err := s.taskService.GetTaskWith(taskId)
	if err != nil {
		logger.Error(err)
	}

	pos, err := s.getPositionWith(taskId)
	if err != nil {
		logger.Error(err)
		// return nil, err
	}

	for mcap := range marketCapChan {
		positionMessage := s.generatePositionMessage(pos, mcap)
		fmt.Println(positionMessage)
		s.publisher.PublishPositionUpdate(task, positionMessage)
		//publish the profit, mcap, position
	}
}

func (s *Service) generatePositionMessage(pos *position.Position, marketCap float64) position.PositionMessage {
	//we need unrealized profit -> remaning tokens (in sol)
	totalPnl, unrealizedPnl := s.getProfitValues(pos, marketCap)

	posMessage := position.PositionMessage{
		BuyTaskId:        pos.PositionId,
		UnrealizedProfit: unrealizedPnl.Text('f', 9),
		RealizedProfit:   pos.FinalizedProfit.Text('f', 9),
		TotalPnL:         totalPnl.Text('f', 9),
		MarketCap:        marketCap,
		RemainingTokens:  pos.TokenRemaining.String(),
	}

	return posMessage
}

func (s *Service) getProfitValues(pos *position.Position, marketCap float64) (totalPnl *big.Float, unrealized *big.Float) {
	solPrice, err := solana_price.GetSolPrice()
	if err != nil {
		logger.Error(err)
	}

	tokenValue := calculateTokenValueInSol(marketCap, pos.InitialTokenAmount.Uint64(), *solPrice)
	unrealizedPnl := new(big.Float).Sub(tokenValue, pos.RemaningCostBasis)

	totalPnL := new(big.Float).Add(pos.FinalizedProfit, unrealizedPnl)

	return totalPnL, unrealizedPnl
}

func (s *Service) getPositionWith(buyTaskId string) (*position.Position, error) {

	pos, ok := s.positions[buyTaskId]
	if !ok {
		return nil, fmt.Errorf("position not found with id: %s", buyTaskId)
	}

	return pos, nil
}

func calculateTokenValueInSol(marketCapUSD float64, tokensRemaining uint64, solPrice float64) (totalValueSOL *big.Float) {
	marketCapSol := new(big.Float).Quo(big.NewFloat(marketCapUSD), big.NewFloat(solPrice))
	totalSupply := new(big.Float).SetInt64(1000000000)

	pricePerTokenSOL := new(big.Float).Quo(marketCapSol, totalSupply)

	tokensFloat := new(big.Float).SetUint64(tokensRemaining)
	totalValueSOL = new(big.Float).Mul(pricePerTokenSOL, tokensFloat)

	return totalValueSOL
}
