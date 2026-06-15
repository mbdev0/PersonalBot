package pool

import (
	"context"
	"math/big"
	"personal_bot/backend/infrastructure/solana_price"
	"strconv"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/sync/errgroup"
)

type PoolBalances struct {
	TokenPoolBalance float64
	WsolPoolBalance  float64
}

func GetTokenBalances(ctx context.Context, poolAtaAddress, wsolAtaAddress string, httpClient *rpc.Client) (PoolBalances, error) {
	var tokenBalance, wsolBalance float64

	group, ctx := errgroup.WithContext(ctx)
	group.Go(
		func() error {
			poolAtaAddy, err := solana.PublicKeyFromBase58(poolAtaAddress)
			if err != nil {
				return err
			}

			res, err := httpClient.GetTokenAccountBalance(ctx, poolAtaAddy, rpc.CommitmentConfirmed)
			if err != nil {
				return err
			}

			poolAtaBalance, err := strconv.ParseFloat(res.Value.UiAmountString, 64)
			if err != nil {
				return err
			}

			tokenBalance = poolAtaBalance
			return nil
		},
	)

	group.Go(
		func() error {
			wsolAtaAddy, err := solana.PublicKeyFromBase58(wsolAtaAddress)
			if err != nil {
				return err
			}

			res, err := httpClient.GetTokenAccountBalance(ctx, wsolAtaAddy, rpc.CommitmentConfirmed)
			if err != nil {
				return err
			}

			wsolAtaBalance, err := strconv.ParseFloat(res.Value.UiAmountString, 32)
			if err != nil {
				return err
			}

			wsolBalance = wsolAtaBalance
			return nil
		},
	)

	if err := group.Wait(); err != nil {
		return PoolBalances{}, err
	}

	return PoolBalances{TokenPoolBalance: tokenBalance, WsolPoolBalance: wsolBalance}, nil
}

func GetMarketCapUSD(poolBalances PoolBalances) (*big.Float, error) {
	wsolPool := big.NewFloat(poolBalances.WsolPoolBalance)
	tokenPool := big.NewFloat(poolBalances.TokenPoolBalance)

	totalSupply := new(big.Float).SetInt64(1_000_000_000)

	pricePerToken := new(big.Float).Quo(wsolPool, tokenPool)

	marketCapSol := new(big.Float).Mul(pricePerToken, totalSupply)

	solPrice, err := solana_price.GetSolPrice()
	if err != nil {
		return nil, err
	}

	solPriceBig := big.NewFloat(*solPrice)
	marketCapUSD := new(big.Float).Mul(solPriceBig, marketCapSol)

	return marketCapUSD, nil
}
