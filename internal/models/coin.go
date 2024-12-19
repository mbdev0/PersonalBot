package models

type DecodedInstruction struct {
	Name     string
	Symbol   string
	IPFS_URL string
}

type MintData struct {
	Signature        string
	Name             string
	Symbol           string
	IPFS_URL         string
	TokenAddr        string
	CreatorAddr      string
	DevHoldingAmount float64
}

type IPFS struct {
	Image    string `json:"image"`
	Twitter  string `json:"twitter"`
	Telegram string `json:"telegram"`
	Website  string `json:"website"`
}

// Compositions of the two structs above
type Coin struct {
	CoinData MintData
	IPFSData IPFS
}
