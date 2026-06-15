package pool

import (
	"personal_bot/backend/pkg/logger"

	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

// given a number of bytes, we borsch decode the pool and then get values off of it
// then we can give this to other parts of pumpfun amm stuff
type Pool struct {
	PoolBump              uint8
	Index                 uint16
	Creator               solana.PK
	BaseMint              solana.PK
	QuoteMint             solana.PK
	LpMint                solana.PK
	PoolBaseTokenAccount  solana.PK
	PoolQuoteTokenAccount solana.PK
	LpSupply              uint64
	CoinCreator           solana.PK
	IsMayhemMode          bool
	IsCashbackCoin        bool
}

func DecodePoolData(data []byte) (Pool, error) {
	var pool Pool
	err := borsh.Deserialize(&pool, data[8:])
	if err != nil {
		return Pool{}, err
	}

	logger.Information("pool data: ", pool)
	return pool, nil
}
