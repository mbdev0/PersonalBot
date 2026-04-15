package position

import (
	"context"
	"fmt"
	"math/big"
	"personal_bot/internal/core/position"
	position_hub "personal_bot/internal/services/subscription_hub/position"
	"sync"

	"github.com/gagliardetto/solana-go"
)

type Service struct {
	positions map[int64]*position.Position
	mu        *sync.Mutex
	subhub    *position_hub.SubscriptionHub
}

func NewPositionService(subhub *position_hub.SubscriptionHub) *Service {
	return &Service{
		positions: map[int64]*position.Position{},
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

func (s *Service) ReportBuy(ctx context.Context, buyReportPayload position.ReportBuyPayload) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entryPrice := new(big.Float).Quo(buyReportPayload.SolSpent, buyReportPayload.TokenAmount)

	newPosition := position.Position{
		PositionId:           buyReportPayload.BuyTaskId,
		TokenAddress:         buyReportPayload.TokenAddress,
		WalletAddress:        buyReportPayload.WalletAddress,
		InitialTokenAmount:   buyReportPayload.TokenAmount,
		TokenRemaining:       buyReportPayload.TokenAmount,
		RemainingCostBasis:   buyReportPayload.SolSpent,
		FinalizedProfit:      big.NewFloat(0),
		InitialSolanaAmount:  buyReportPayload.SolSpent,
		EntryPrice:           entryPrice,
		MarketCapEntry:       buyReportPayload.MarketCap,
		TotalMarketCapExit:   big.NewFloat(0),
		NumberOfSells:        0,
		AverageMarketCapExit: big.NewFloat(0),
	}

	s.positions[newPosition.PositionId] = &newPosition

	//publish buy to positionhub -> handle ctx in there
	return s.subhub.PublishPositionCreate(ctx, &newPosition)
}

func (s *Service) ReportSell(buyTaskId int64, tokensSold *big.Float, solRecieved *big.Float, marketCapSale *big.Float) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pos, ok := s.positions[buyTaskId]
	if !ok {
		return fmt.Errorf("position not found with id: %d", buyTaskId)
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
	pos.NumberOfSells++
	pos.TotalMarketCapExit = new(big.Float).Add(pos.TotalMarketCapExit, marketCapSale)
	pos.AverageMarketCapExit = new(big.Float).Quo(pos.TotalMarketCapExit, new(big.Float).SetInt64(int64(pos.NumberOfSells)))
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

func (s *Service) GetById(id int64) (*position.Position, error) {
	pos, ok := s.positions[id]
	if !ok {
		return nil, fmt.Errorf("position not found with id %d", id)
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

func (s *Service) Subscribe(id int64, isInternalSub bool) (*position_hub.Subscription, error) {
	sub, err := s.subhub.Subscribe(id, isInternalSub, nil)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *Service) Unsubscribe(id int64, isInternalSub bool) error {
	err := s.subhub.Unsubscribe(id, isInternalSub)
	if err != nil {
		return err
	}

	return nil
}
