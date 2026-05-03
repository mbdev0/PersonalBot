package models

import "github.com/gagliardetto/solana-go"

type DecodedCreateInstructionV1 struct {
	Name    string
	Symbol  string
	IpfsUrl string
}

type DecodedCreateInstructionV2 struct {
	Name         string
	Symbol       string
	IpfsUrl      string
	Creator      solana.PublicKey
	IsMayhemMode bool
}

type DecodedCreateInstruction struct {
	Name         string
	Symbol       string
	IpfsUrl      string
	Creator      solana.PublicKey
	IsMayhemMode bool
}

type Coin struct {
	CoinData MintData
	IPFSData IPFS
}

type MintData struct {
	Signature        string
	Name             string
	Symbol           string
	IpfsUrl          string
	TokenAddr        string
	CreatorAddr      string
	BondingCurveAddr string
	DevHoldingAmount float64
}

type IPFS struct {
	ImageURL    string `json:"image"`
	TwitterURL  string `json:"twitter"`
	TelegramURL string `json:"telegram"`
	WebsiteURL  string `json:"website"`
}
