package mapper

import (
	"fmt"
	"math/big"
	"personal_bot/infrastructure/persistence/models"
	"personal_bot/internal/core/position"

	"github.com/gagliardetto/solana-go"
)

func MapPositionToRepo(pos position.Position) (repo models.PositionRepository, err error) {
	if pos.InitialTokenAmount == nil {
		return repo, fmt.Errorf("InitialTokenAmount was nil")
	}
	initalTokenAmount := pos.InitialTokenAmount.String()

	if pos.TokenRemaining == nil {
		return repo, fmt.Errorf("TokenRemaining was nil")
	}
	tokenRemaining := pos.TokenRemaining.String()

	if pos.RemainingCostBasis == nil {
		return repo, fmt.Errorf("RemainingCostBasis was nil")
	}
	remainingCostBasis := pos.RemainingCostBasis.String()

	if pos.FinalizedProfit == nil {
		return repo, fmt.Errorf("FinalizedProfit was nil")
	}
	finalizedProfit := pos.FinalizedProfit.String()

	if pos.InitialSolanaAmount == nil {
		return repo, fmt.Errorf("InitialSolanaAmount was nil")
	}
	initialSolanaAmount := pos.InitialSolanaAmount.String()

	if pos.EntryPrice == nil {
		return repo, fmt.Errorf("EntryPrice was nil")
	}
	entryPrice := pos.EntryPrice.String()

	if pos.MarketCapEntry == nil {
		return repo, fmt.Errorf("MarketCapEntry was nil")
	}
	mcapEntry := pos.MarketCapEntry.String()

	if pos.AverageMarketCapExit == nil {
		return repo, fmt.Errorf("AverageMarketCapExit was nil")
	}
	averageMcapExit := pos.AverageMarketCapExit.String()

	if pos.TotalMarketCapExit == nil {
		return repo, fmt.Errorf("TotalMarketCapExit was nil")
	}
	totalMcapExit := pos.TotalMarketCapExit.String()

	repo = models.PositionRepository{
		PositionId:         pos.PositionId,
		StrategyId:         pos.StrategyId,
		TokenAddress:       pos.TokenAddress.String(),
		WalletAddress:      pos.WalletAddress.String(),
		InitialTokenAmount: initalTokenAmount,
		TokensRemaining:    tokenRemaining,
		RemainingCostBasis: remainingCostBasis,
		FinalizedProfit:    finalizedProfit,
		InitialAmount:      initialSolanaAmount,
		EntryPrice:         entryPrice,
		McapEntry:          mcapEntry,
		AverageMcapExit:    averageMcapExit,
		TotalMarketCapExit: totalMcapExit,
		NumberOfSells:      int64(pos.NumberOfSells),
		AddressForUrl:      pos.AddressForUrl,
	}

	return repo, nil

}

func MapRepoToPosition(repo models.PositionRepository) (pos position.Position, err error) {
	initalTokenAmount, success := new(big.Float).SetString(repo.InitialTokenAmount)
	if !success {
		return pos, fmt.Errorf("not able to cast initalTokenAmount to big float")
	}

	tokenRemaining, success := new(big.Float).SetString(repo.TokensRemaining)
	if !success {
		return pos, fmt.Errorf("not able to cast tokenRemaining to big float")
	}

	remainingCostBasis, success := new(big.Float).SetString(repo.RemainingCostBasis)
	if !success {
		return pos, fmt.Errorf("not able to cast remainingCostBasis to big float")
	}

	finalizedProfit, success := new(big.Float).SetString(repo.FinalizedProfit)
	if !success {
		return pos, fmt.Errorf("not able to cast finalizedProfit to big float")
	}

	initialAmount, success := new(big.Float).SetString(repo.InitialAmount)
	if !success {
		return pos, fmt.Errorf("not able to cast initialAmount to big float")
	}

	entryPrice, success := new(big.Float).SetString(repo.EntryPrice)
	if !success {
		return pos, fmt.Errorf("not able to cast entryPrice to big float")
	}

	mcapEntry, success := new(big.Float).SetString(repo.McapEntry)
	if !success {
		return pos, fmt.Errorf("not able to cast mcapEntry to big float")
	}

	averageMcapExit, success := new(big.Float).SetString(repo.AverageMcapExit)
	if !success {
		return pos, fmt.Errorf("not able to cast averageMcapExit to big float")
	}

	totalMcapExit, success := new(big.Float).SetString(repo.TotalMarketCapExit)
	if !success {
		return pos, fmt.Errorf("not able to cast totalMcapExit to big float")
	}

	return position.Position{
		PositionId:           repo.PositionId,
		StrategyId:           repo.StrategyId,
		TokenAddress:         solana.MustPublicKeyFromBase58(repo.TokenAddress),
		WalletAddress:        solana.MustPublicKeyFromBase58(repo.WalletAddress),
		InitialTokenAmount:   initalTokenAmount,
		TokenRemaining:       tokenRemaining,
		RemainingCostBasis:   remainingCostBasis,
		FinalizedProfit:      finalizedProfit,
		InitialSolanaAmount:  initialAmount,
		EntryPrice:           entryPrice,
		MarketCapEntry:       mcapEntry,
		AverageMarketCapExit: averageMcapExit,
		TotalMarketCapExit:   totalMcapExit,
		NumberOfSells:        int(repo.NumberOfSells),
		AddressForUrl:        repo.AddressForUrl,
	}, nil
}
