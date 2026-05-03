package models

import "github.com/gagliardetto/solana-go"

type TradeEvent struct {
	Mint                  solana.PublicKey
	SolAmount             uint64
	TokenAmount           uint64
	IsBuy                 bool
	User                  solana.PublicKey
	Timestamp             int64
	VirtualSolReserves    uint64
	VirtualTokenReserves  uint64
	RealSolReserves       uint64
	RealTokenReserves     uint64
	FeeRecipient          solana.PublicKey
	FeeBasisPoints        uint64
	Fee                   uint64
	Creator               solana.PublicKey
	CreatorFeeBasisPoints uint64
	CreatorFee            uint64
	TrackVolume           bool
	TotalUnclaimedTokens  uint64
	TotalClaimedTokens    uint64
	CurrentSolVolume      uint64
	LastUpdateTimestamp   int64
	IxName                string
}
