package strategies

import (
	"fmt"
	"math/big"
	rpcgroups "personal_bot/internal/core/rpc_groups"
	"personal_bot/internal/core/wallets"

	"github.com/gagliardetto/solana-go"
)

type SpamPatch struct {
	Program          *string
	BuyAmount        *big.Int
	BuyFee           *float64
	Token            *solana.PublicKey
	Slippage         *float64
	ComputeUnits     *float64
	Wallet           *wallets.SolanaWallet
	SellFee          *float64
	RPCGroup         *rpcgroups.RPCGroup
	Retries          *uint16
	RetriesDelayMS   *uint32
	NumberOfSubTasks *uint32
	StartTime        *uint64
}

func (sp *SpamPatch) ApplyTo(task Task) error {
	spam, ok := task.(*Spam)
	if !ok {
		return fmt.Errorf("unable to patch task - was given a task that was not spam")
	}

	if sp.Program != nil {
		spam.Program = *sp.Program
	}

	if sp.BuyAmount != nil {
		spam.BuyAmount = sp.BuyAmount
	}

	if sp.BuyFee != nil {
		spam.BuyFee = *sp.BuyFee
	}

	if sp.ComputeUnits != nil {
		spam.ComputeUnits = *sp.ComputeUnits
	}

	if sp.Wallet != nil {
		spam.Wallet = *sp.Wallet
	}

	if sp.SellFee != nil {
		spam.SellFee = sp.SellFee
	}

	if sp.RPCGroup != nil {
		spam.RPCGroup = *sp.RPCGroup
	}
	if sp.Retries != nil {
		spam.Retries = *sp.Retries
	}

	if sp.RetriesDelayMS != nil {
		spam.RetriesDelayMS = *sp.RetriesDelayMS
	}
	if sp.NumberOfSubTasks != nil {
		spam.NumberOfSubTasks = *sp.NumberOfSubTasks
	}
	if sp.StartTime != nil {
		spam.StartTime = *sp.StartTime
	}

	if sp.Token != nil {
		spam.Token = *sp.Token
	}

	return nil
}
