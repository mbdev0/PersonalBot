package position

import (
	"fmt"
	"math/big"
	"pump_fun/internal/core/position"
	position_hub "pump_fun/internal/services/subscription_hub/position"
	"sync"

	"github.com/gagliardetto/solana-go"
)

type Service struct {
	positions map[string]*position.Position
	mu        *sync.Mutex
	subhub    *position_hub.SubscriptionHub
}

func NewPositionService(subhub *position_hub.SubscriptionHub) Service {
	return Service{
		positions: map[string]*position.Position{},
		mu:        &sync.Mutex{},
		subhub:    subhub,
	}
}

func (s *Service) FindPositionIfExists(token solana.PublicKey, walletAddress solana.PublicKey) (*position.Position, bool) {
	for _, value := range s.positions {
		if value.TokenAddress == token && value.WalletAddress == walletAddress {
			return value, true
		}
	}

	return nil, false
}

func (s *Service) ReportBuy(buytaskid string, tokenaddress solana.PublicKey, walletAddress solana.PublicKey, tokenAmount *big.Float, solSpent *big.Float) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newPosition := position.Position{
		PositionId:         buytaskid,
		TokenAddress:       tokenaddress,
		WalletAddress:      walletAddress,
		InitialTokenAmount: tokenAmount,
		TokenRemaining:     tokenAmount,
		RemainingCostBasis: solSpent,
		FinalizedProfit:    big.NewFloat(0),
	}

	s.positions[newPosition.PositionId] = &newPosition

	//publish buy to positionhub -> handle ctx in there
	s.subhub.PublishPositionCreate(&newPosition)
}

func (s *Service) ReportSell(buyTaskId string, tokensSold *big.Float, solRecieved *big.Float) error {
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
		tokensSold,
		pos.TokenRemaining,
	)

	costBasisofSoldTokens := new(big.Float).Mul(sellRatio, pos.RemainingCostBasis)
	realizedBasisFromSale := new(big.Float).Sub(solRecieved, costBasisofSoldTokens)

	pos.RemainingCostBasis = new(big.Float).Sub(pos.RemainingCostBasis, costBasisofSoldTokens)
	pos.FinalizedProfit = new(big.Float).Add(pos.FinalizedProfit, realizedBasisFromSale)
	pos.TokenRemaining = new(big.Float).Sub(pos.TokenRemaining, tokensSold)

	//publish sell
	err := s.subhub.PublishPositionUpdate(pos)
	if err != nil {
		return err
	}

	// if pos.remainingTokens <= 0 -> publish -> close stream
	isTokensRemaining := pos.TokenRemaining.Cmp(big.NewFloat(0)) == 1
	if !isTokensRemaining {
		err := s.subhub.PublishPositionStop(pos)
		if err != nil {
			return err
		}
	}

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

func (s *Service) Subscribe(id string) (*position_hub.Subscription, error) {
	sub, err := s.subhub.Subscribe(id)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *Service) Unsubscribe(id string) error {
	err := s.subhub.Unsubscribe(id)
	if err != nil {
		return err
	}

	return nil
}

// Subscribe
//	handler
//	controller
//	pos_service
//	subhub ->
//	streamer -> streamer needs to publish messages...
// 			 -> passes position messages to channel given from subhub
// subhub 	 -> reads each message and publishes

// Unsubscribe
//	handler
//	controller
//	pos_service
//	subhub
//	streamer

// Unsub via no tokens remaining
// pos_Service
// streamer -> stop
