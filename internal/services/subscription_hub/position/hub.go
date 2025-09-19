package position

import (
	"context"
	"fmt"
	"math/big"
	"pump_fun/infrastructure/solana_price"
	"pump_fun/internal/core/position"
	"pump_fun/internal/monitoring"
	"pump_fun/internal/solana/programs/pumpfun/pda"
	"pump_fun/pkg/logger"
	"sync"
)

type Subscription struct {
	sub_id    string
	SubChan   chan position.PositionMessage
	cancel    func()
	cancelCtx context.CancelFunc
}

type SubscriptionHub struct {
	subscriptions map[string]*Subscription
	last          map[string]*position.PositionMessage
	bufferSize    int
	mu            *sync.Mutex
}

func NewSubscriptionHub() *SubscriptionHub {
	return &SubscriptionHub{
		subscriptions: map[string]*Subscription{},
		last:          map[string]*position.PositionMessage{},
		bufferSize:    1000,
		mu:            &sync.Mutex{},
	}
}

func (sh *SubscriptionHub) Subscribe(positionId string) (*Subscription, error) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if _, ok := sh.subscriptions[positionId]; ok {
		return nil, fmt.Errorf("position already subbed to")
	}

	sub := &Subscription{
		sub_id:  positionId,
		SubChan: make(chan position.PositionMessage, sh.bufferSize),
		cancel:  sh.cancel(positionId),
	}

	if last, ok := sh.last[positionId]; ok {
		select {
		case sub.SubChan <- *last:
		default:
			logger.Error("not able to push message onto channel - id: ", sub.sub_id)
		}
	}

	sh.subscriptions[positionId] = sub
	return sub, nil
}

func (sh *SubscriptionHub) cancel(positionId string) func() {
	return func() {
		sub := sh.subscriptions[positionId]
		close(sub.SubChan)
		sub.cancelCtx()
		delete(sh.subscriptions, positionId)
		delete(sh.last, positionId)
	}
}

func (sh *SubscriptionHub) Unsubscribe(positionId string) error {
	//get the sub if exists
	pos, ok := sh.subscriptions[positionId]
	if !ok {
		return fmt.Errorf("not able to find positon with id: %s", positionId)
	}
	// run the cancel func
	sh.mu.Lock()
	pos.cancel()
	sh.mu.Unlock()
	return nil
}

func (sh *SubscriptionHub) publish(id string, posMessage *position.PositionMessage) {
	//if the sub exists push to chan
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sub, ok := sh.subscriptions[id]; ok {
		sub.SubChan <- *posMessage
	}

	if _, ok := sh.last[id]; ok {
		sh.last[id] = posMessage
	}
}

func (sh *SubscriptionHub) PublishPositionUpdate(pos *position.Position) error {
	//	get position's last message
	posMessage, ok := sh.last[pos.PositionId]
	if !ok {
		return fmt.Errorf("position not found")
	}
	//	get mcap
	sh.mu.Lock()
	mcap, _, err := big.ParseFloat(posMessage.MarketCap, 10, 9, big.ToNearestEven)
	if err != nil {
		return err
	}
	//	regenerate profit etc.
	message := sh.generatePositionMessage(pos, mcap)
	sh.mu.Unlock()

	//	publish
	sh.publish(pos.PositionId, &message)
	return nil
}

func (sh *SubscriptionHub) PublishPositionCreate(p *position.Position) {
	ctx, cancel := context.WithCancel(context.Background())

	sub := Subscription{
		sub_id:    p.PositionId,
		SubChan:   make(chan position.PositionMessage, sh.bufferSize),
		cancel:    sh.cancel(p.PositionId),
		cancelCtx: cancel,
	}
	sh.mu.Lock()
	sh.subscriptions[sub.sub_id] = &sub
	sh.mu.Unlock()

	posMessage := position.PositionMessage{
		MessageType:      position.Created,
		BuyTaskId:        sub.sub_id,
		UnrealizedProfit: "",
		RealizedProfit:   p.FinalizedProfit.Text('f', 9),
		RemainingTokens:  p.TokenRemaining.Text('f', 9),
		Message:          "Position Created",
	}

	sh.publish(p.PositionId, &posMessage)

	go func(s *Subscription, c context.Context, pos *position.Position) {
		marketCapChan := make(chan *big.Float, sh.bufferSize)

		bondingCurveAddress, err := pda.GetBondingCurveAddress(pos.TokenAddress.String())
		if err != nil {
			logger.Error(err)
		}
		// start marketcap streaming
		go monitoring.StartMarketCapMonitor(ctx, bondingCurveAddress, marketCapChan)

		// then foreach in mcap
		for mcap := range marketCapChan {
			positionMessage := sh.generatePositionMessage(pos, mcap)
			positionMessage.MessageType = position.Update
			fmt.Println(positionMessage)
			//publish the profit, mcap, position
			sh.publish(s.sub_id, &positionMessage)
		}
	}(&sub, ctx, p)
}

func (sh *SubscriptionHub) PublishPositionStop(pos *position.Position) error {

	posMessage, ok := sh.last[pos.PositionId]
	if !ok {
		return fmt.Errorf("could not find an open position: %s", pos.PositionId)
	}

	message := position.PositionMessage{
		MessageType:      position.Stopped,
		BuyTaskId:        pos.PositionId,
		UnrealizedProfit: posMessage.UnrealizedProfit,
		RealizedProfit:   posMessage.RealizedProfit,
		TotalPnL:         posMessage.TotalPnL,
		MarketCap:        posMessage.MarketCap,
		RemainingTokens:  posMessage.RemainingTokens,
		Message:          "Position Closed",
	}
	//when the position is closed

	sh.publish(pos.PositionId, &message)
	return nil
}

func (sh *SubscriptionHub) generatePositionMessage(pos *position.Position, marketCap *big.Float) position.PositionMessage {
	//we need unrealized profit -> remaning tokens (in sol)
	totalPnl, unrealizedPnl := sh.getProfitValues(pos, marketCap)

	posMessage := position.PositionMessage{
		BuyTaskId:        pos.PositionId,
		UnrealizedProfit: unrealizedPnl.Text('f', 9),
		RealizedProfit:   pos.FinalizedProfit.Text('f', 9),
		TotalPnL:         totalPnl.Text('f', 9),
		MarketCap:        marketCap.String(),
		RemainingTokens:  pos.TokenRemaining.String(),
	}

	return posMessage
}

func (sh *SubscriptionHub) getProfitValues(pos *position.Position, marketCap *big.Float) (totalPnl *big.Float, unrealized *big.Float) {
	solPrice, err := solana_price.GetSolPrice()
	if err != nil {
		logger.Error(err)
	}

	tokenValue := sh.calculateTokenValueInSol(marketCap, pos.InitialTokenAmount, *solPrice)
	unrealizedPnl := new(big.Float).Sub(tokenValue, pos.RemainingCostBasis)

	totalPnL := new(big.Float).Add(pos.FinalizedProfit, unrealizedPnl)

	return totalPnL, unrealizedPnl
}

func (sh *SubscriptionHub) calculateTokenValueInSol(marketCapUSD *big.Float, tokensRemaining *big.Float, solPrice float64) (totalValueSOL *big.Float) {
	marketCapSol := new(big.Float).Quo(marketCapUSD, big.NewFloat(solPrice))
	totalSupply := new(big.Float).SetInt64(1000000000)

	pricePerTokenSOL := new(big.Float).Quo(marketCapSol, totalSupply)

	totalValueSOL = new(big.Float).Mul(pricePerTokenSOL, tokensRemaining)

	return totalValueSOL
}
