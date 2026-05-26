package position

import (
	"context"
	"fmt"
	"math/big"
	"personal_bot/infrastructure/solana_price"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/core/position"
	"personal_bot/internal/core/program"
	rpcgroups "personal_bot/internal/core/rpc_groups"
	"personal_bot/internal/services/settings"
	datastructures "personal_bot/pkg/data_structures"
	"personal_bot/pkg/logger"
	"sync"

	"github.com/gagliardetto/solana-go/rpc"
)

type Subscription struct {
	Sub_id         int64
	SubChan        chan position.PositionMessage
	cancel         func()
	cancelCtx      context.CancelFunc
	hasInternalSub bool
	hasClientSub   bool
}

type SubscriptionHub struct {
	subscriptions   map[int64]*Subscription
	activePositions datastructures.Map[int64, *position.Position]
	last            map[int64]*position.PositionMessage
	bufferSize      int
	mu              *sync.Mutex
	settings        *settings.Service
}

func NewSubscriptionHub(settings *settings.Service) *SubscriptionHub {
	return &SubscriptionHub{
		subscriptions:   map[int64]*Subscription{},
		activePositions: datastructures.NewMap[int64, *position.Position](),
		last:            map[int64]*position.PositionMessage{},
		bufferSize:      1000,
		mu:              &sync.Mutex{},
		settings:        settings,
	}
}

func (sh *SubscriptionHub) Subscribe(positionId int64, isInternalSub bool, rpcGroup *rpcgroups.RPCNode) (*Subscription, error) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if _, ok := sh.subscriptions[positionId]; ok {
		return nil, fmt.Errorf("position already subbed to")
	}

	ctx, cancel := context.WithCancel(context.Background())

	sub := &Subscription{
		Sub_id:    positionId,
		SubChan:   make(chan position.PositionMessage, sh.bufferSize),
		cancel:    sh.cancel(positionId),
		cancelCtx: cancel,
	}

	if isInternalSub {
		sub.hasInternalSub = true
	} else {
		sub.hasClientSub = true
	}

	p, ok := sh.activePositions.Get(positionId)
	if !ok {
		return nil, fmt.Errorf("position has not been created yet")
	}

	go func(c context.Context, s *Subscription, pos *position.Position) {

		// marketCapChan := make(chan *big.Float, sh.bufferSize)

		node, err := sh.resolveRPCNode(c, rpcGroup)
		if err != nil {
			logger.Error(err)
			return
		}

		prog := program.Resolve(pos.Program)
		stream := prog.NewMarketCapStream(c, pos.MonitoringAddress, node)

		// go monitoring.StartMarketCapMonitor(ctx, pos.MonitoringAddress, marketCapChan, node)

		// start marketcap streaming
		// then foreach in mcap
		for mcap := range stream {
			positionMessage := sh.generatePositionMessage(pos, mcap)
			positionMessage.MessageType = position.Update
			//publish the profit, mcap, position
			sh.publish(s.Sub_id, &positionMessage)
		}
	}(ctx, sub, p)

	if last, ok := sh.last[positionId]; ok {
		select {
		case sub.SubChan <- *last:
		default:
			logger.Error("not able to push message onto channel - id: ", sub.Sub_id)
		}
	}

	sh.subscriptions[positionId] = sub
	return sub, nil
}

func (sh *SubscriptionHub) resolveRPCNode(c context.Context, rpcGroup *rpcgroups.RPCNode) (rpcgroups.RPCNode, error) {
	if rpcGroup != nil {
		return *rpcGroup, nil
	}

	settings, err := sh.settings.GetSettings(c)
	if err != nil {
		return rpcgroups.RPCNode{}, err
	}

	return rpcgroups.RPCNode{
		Http: rpc.New(settings.PositionNodes.HTTPNode),
		WS:   settings.PositionNodes.WSNode,
	}, nil
}

func (sh *SubscriptionHub) cancel(positionId int64) func() {
	return func() {
		sub := sh.subscriptions[positionId]
		sub.cancelCtx()
		// close(sub.SubChan)
		delete(sh.subscriptions, positionId)
		delete(sh.last, positionId)
	}
}

func (sh *SubscriptionHub) Unsubscribe(positionId int64, isInternalSub bool) {
	//get the sub if exists
	pos, ok := sh.subscriptions[positionId]
	if !ok {
		return
	}

	if isInternalSub {
		pos.hasInternalSub = false
	} else {
		pos.hasClientSub = false
	}
	// run the cancel func
	sh.mu.Lock()
	//if there are no subscribers on the task
	if !pos.hasClientSub && !pos.hasInternalSub {
		pos.cancel()
	}
	sh.mu.Unlock()
}

func (sh *SubscriptionHub) WaitForCreate(id int64) (*position.Position, bool) {
	pos, ok := sh.activePositions.WaitForEntry(id)
	return pos, ok
	//wait until position is created in activePositions

}

func (sh *SubscriptionHub) publish(id int64, posMessage *position.PositionMessage) {
	//if the sub exists push to chan
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sub, ok := sh.subscriptions[id]; ok {
		sub.SubChan <- *posMessage
	}

	sh.last[id] = posMessage
}

func (sh *SubscriptionHub) PublishPositionUpdate(pos *position.Position) error {
	//	get position's last message
	posMessage, ok := sh.last[pos.PositionId]
	if !ok {
		return fmt.Errorf("position not found")
	}
	//	get mcap
	sh.mu.Lock()
	mcap := posMessage.MarketCap
	if mcap == nil {
		return fmt.Errorf("posMessage.MarketCap was nil")
	}

	//	regenerate profit etc.
	message := sh.generatePositionMessage(pos, *mcap)
	message.Message = "Position Update"
	message.MessageType = position.Update
	sh.mu.Unlock()

	//	publish
	sh.publish(pos.PositionId, &message)
	return nil
}

func (sh *SubscriptionHub) PublishPositionCreate(ctx context.Context, p *position.Position) error {
	sh.mu.Lock()
	sh.activePositions.Set(p.PositionId, p)
	sh.mu.Unlock()

	if p.MarketCapEntry == nil {
		return fmt.Errorf("MarketCapEntry was nil")
	}

	posMsg := sh.generatePositionMessage(p, *p.MarketCapEntry)
	posMsg.MessageType = position.Created
	posMsg.Message = "Position Created"

	sh.publish(p.PositionId, &posMsg)
	return nil
}

func (sh *SubscriptionHub) PublishPositionStop(pos *position.Position) error {

	posMessage, ok := sh.last[pos.PositionId]
	if !ok {
		return fmt.Errorf("could not find an open position: %d", pos.PositionId)
	}

	message := position.PositionMessage{
		MessageType:         position.Stopped,
		StrategyId:          pos.StrategyId,
		BuyTaskId:           pos.PositionId,
		UnrealizedProfit:    posMessage.UnrealizedProfit,
		RealizedProfit:      posMessage.RealizedProfit,
		TotalPnL:            posMessage.TotalPnL,
		MarketCap:           posMessage.MarketCap,
		RemainingTokens:     posMessage.RemainingTokens,
		InitialSolanaAmount: posMessage.InitialSolanaAmount,
		Message:             "Position Closed",
	}
	//when the position is closed
	sh.mu.Lock()
	sh.activePositions.Delete(pos.PositionId)
	sh.mu.Unlock()

	sh.publish(pos.PositionId, &message)
	return nil
}

func (sh *SubscriptionHub) generatePositionMessage(pos *position.Position, mcap big.Float) position.PositionMessage {
	//we need unrealized profit -> remaning tokens (in sol)
	marketCap := big.NewFloat(0).Set(&mcap)
	totalPnl, unrealizedPnl, currentPrice := sh.getProfitValues(pos, marketCap)

	totalPnl.Quo(totalPnl, big.NewFloat(constants.LamportsConversion))
	unrealizedPnl.Quo(unrealizedPnl, big.NewFloat(constants.LamportsConversion))
	finalizedProfit := new(big.Float).Quo(pos.FinalizedProfit, big.NewFloat(constants.LamportsConversion))
	tokensRemaining := new(big.Float).Quo(pos.TokenRemaining, big.NewFloat(constants.TokenAmountDecimals))

	posMessage := position.PositionMessage{
		BuyTaskId:           pos.PositionId,
		StrategyId:          pos.StrategyId,
		UnrealizedProfit:    unrealizedPnl,
		RealizedProfit:      finalizedProfit,
		TotalPnL:            totalPnl,
		MarketCap:           marketCap,
		RemainingTokens:     tokensRemaining,
		InitialSolanaAmount: pos.InitialSolanaAmount,
		CurrentPrice:        currentPrice,
	}

	return posMessage
}

func (sh *SubscriptionHub) getProfitValues(pos *position.Position, marketCap *big.Float) (totalPnl *big.Float, unrealized *big.Float, currentPrice *big.Float) {
	solPrice, err := solana_price.GetSolPrice()
	if err != nil {
		logger.Error(err)
	}

	tokenValue, currentPrice := sh.calculateTokenValueAndPrice(marketCap, pos.TokenRemaining, *solPrice)
	unrealizedPnl := new(big.Float).Sub(tokenValue, pos.RemainingCostBasis)
	totalPnL := new(big.Float).Add(pos.FinalizedProfit, unrealizedPnl)

	return totalPnL, unrealizedPnl, currentPrice
}

func (sh *SubscriptionHub) calculateTokenValueAndPrice(marketCapUSD *big.Float, tokensRemaining *big.Float, solPrice float64) (totalValueSOL *big.Float, pricePerToken *big.Float) {
	//should solve the nil ref?
	if marketCapUSD == nil {
		marketCapUSD.SetFloat64(0)
	}
	marketCapSol := new(big.Float).Quo(marketCapUSD, big.NewFloat(solPrice))

	//total supply -> get data from
	totalSupply := new(big.Float).SetInt64(1000000000)

	pricePerTokenSOL := new(big.Float).Quo(marketCapSol, totalSupply)

	tokensRemainingWhole := new(big.Float).Quo(tokensRemaining, big.NewFloat(constants.TokenAmountDecimals))

	totalValueSOL = new(big.Float).Mul(pricePerTokenSOL, tokensRemainingWhole)

	return totalValueSOL, pricePerTokenSOL
}
