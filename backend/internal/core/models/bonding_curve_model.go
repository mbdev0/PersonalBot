package models

import (
	"github.com/gagliardetto/solana-go"
)

type BondingCurve struct {
	VirtualTokenReserves uint64
	VirtualSolReserves   uint64
	RealTokenReserves    uint64
	RealSolReserves      uint64
	TokenTotalSupply     uint64
	Complete             bool
	Creator              *solana.PublicKey
	IsMayhemMode         bool
	IsCashbackCoin       bool
}

type BondingCurveV0 struct {
	VirtualTokenReserves uint64
	VirtualSolReserves   uint64
	RealTokenReserves    uint64
	RealSolReserves      uint64
	TokenTotalSupply     uint64
	Complete             bool
}

type BondingCurveV1 struct {
	VirtualTokenReserves uint64
	VirtualSolReserves   uint64
	RealTokenReserves    uint64
	RealSolReserves      uint64
	TokenTotalSupply     uint64
	Complete             bool
	Creator              solana.PublicKey
	IsMayhemMode         bool
}

// TODO: when we care about cashback coins only, expand v2
type BondingCurveV2 struct {
	VirtualTokenReserves uint64
	VirtualSolReserves   uint64
	RealTokenReserves    uint64
	RealSolReserves      uint64
	TokenTotalSupply     uint64
	Complete             bool
	Creator              solana.PublicKey
	IsMayhemMode         bool
	IsCashbackCoin       bool
}
