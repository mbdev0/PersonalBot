package monitoring

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"math"
	"math/big"
	"personal_bot/internal/core/constants"
	datastream "personal_bot/internal/solana/monitoring/data_stream"
	"personal_bot/internal/solana/monitoring/models"
	"personal_bot/internal/solana/programs/pumpfun/pumpfun_amm/pool"
	jsonparse "personal_bot/pkg/json_parse"
	"personal_bot/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type PumpfunAMMMarketCapMonitor struct {
	datastream  datastream.DataStream
	poolAddress string
	ws          string
	httpClient  *rpc.Client
}

func NewPumpfunAMMMarketCapMonitor(datastream datastream.DataStream, poolAddress, ws string, httpClient *rpc.Client) *PumpfunAMMMarketCapMonitor {
	return &PumpfunAMMMarketCapMonitor{
		datastream:  datastream,
		poolAddress: poolAddress,
		ws:          ws,
		httpClient:  httpClient,
	}
}

type PoolPair struct {
	ATokens *float64
	BTokens *float64
}

func (pmm *PumpfunAMMMarketCapMonitor) StreamMarketCap(ctx context.Context) <-chan big.Float {
	poolAddress, err := solana.PublicKeyFromBase58(pmm.poolAddress)
	if err != nil {
		logger.Error(err)
		return nil
	}

	poolBytes, err := pool.GetPoolDataBytes(ctx, poolAddress, pmm.httpClient)
	if err != nil {
		logger.Error(err)
		return nil
	}

	poolData, err := pool.DecodePoolData(poolBytes)
	if err != nil {
		logger.Error(err)
		return nil
	}

	var wsolAccount solana.PK
	var tokenAccount solana.PK

	if poolData.QuoteMint == solana.MustPublicKeyFromBase58(constants.WSOLTokenAddress) {
		wsolAccount = poolData.PoolQuoteTokenAccount
		tokenAccount = poolData.PoolBaseTokenAccount
	} else {
		wsolAccount = poolData.PoolBaseTokenAccount
		tokenAccount = poolData.PoolQuoteTokenAccount
	}

	streamWsol, err := pmm.datastream.SubscribeToAccountStream(pmm.ws, wsolAccount.String())
	if err != nil {
		logger.Error(err)
		return nil
	}
	streamToken, err := pmm.datastream.SubscribeToAccountStream(pmm.ws, tokenAccount.String())
	if err != nil {
		logger.Error(err)
		return nil
	}

	cleanup := func() {
		pmm.datastream.Unsubscribe(pmm.ws, wsolAccount.String())
		pmm.datastream.Unsubscribe(pmm.ws, tokenAccount.String())
	}

	cache := map[uint64]*PoolPair{}

	output := make(chan big.Float)
	go func() {
		defer close(output)
		defer cleanup()

		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-streamWsol:
				if !ok {
					return
				}

				slot, tokenAmount, err := pmm.parse(data)
				if err != nil {
					logger.Error(err)
					continue
				}

				if _, ok := cache[slot]; !ok {
					cache[slot] = &PoolPair{}
				}

				cache[slot].ATokens = &tokenAmount

				if cache[slot].BTokens != nil {
					//then we can gen a marketcap and output
					pmm.output(output, *cache[slot].ATokens, *cache[slot].BTokens)
					delete(cache, slot)
				}

			case data, ok := <-streamToken:
				if !ok {
					return
				}

				slot, tokenAmount, err := pmm.parse(data)
				if err != nil {
					logger.Error(err)
					continue
				}

				if _, ok := cache[slot]; !ok {
					cache[slot] = &PoolPair{}
				}

				cache[slot].BTokens = &tokenAmount

				if cache[slot].ATokens != nil {
					pmm.output(output, *cache[slot].ATokens, *cache[slot].BTokens)
					delete(cache, slot)
				}
			}

			// we should routinely delete cache -> delete all orphaned messages
			pmm.deleteOrphaned(cache)
		}
	}()

	return output
}

const LIMIT = 10

func (pmm *PumpfunAMMMarketCapMonitor) deleteOrphaned(cache map[uint64]*PoolPair) {
	if len(cache) < LIMIT {
		return
	}

	for k, v := range cache {
		if v.ATokens == nil || v.BTokens == nil {
			delete(cache, k)
		}
	}
}

func (pmm *PumpfunAMMMarketCapMonitor) parse(data []byte) (slot uint64, tokenAmount float64, err error) {
	res, err := jsonparse.Decode[models.AccountSubscribeModel](data)
	if err != nil {
		logger.Error(err)
		return
	}

	accountData, err := base64.StdEncoding.DecodeString(res.Params.Result.Value.Data[0])
	if len(accountData) < 72 {
		return
	}

	tokenAmountUint := binary.LittleEndian.Uint64(accountData[64:72])
	tokenAmount = float64(tokenAmountUint)

	return res.Params.Result.Context.Slot, tokenAmount, nil
}

func (pmm *PumpfunAMMMarketCapMonitor) output(out chan big.Float, a, b float64) {
	marketcap, err := pool.GetMarketCapUSD(pool.PoolBalances{
		WsolPoolBalance:  a / math.Pow(10, 9),
		TokenPoolBalance: b / math.Pow(10, 6),
	})

	if err != nil {
		logger.Error(err)
		return
	}

	if marketcap != nil {
		out <- *marketcap
	}

}
