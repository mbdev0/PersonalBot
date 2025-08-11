package models

type DecodedCreateInstruction struct {
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
	BondingCurveAddr string
	DevHoldingAmount float64
}

type IPFS struct {
	ImageURL    string `json:"image"`
	TwitterURL  string `json:"twitter"`
	TelegramURL string `json:"telegram"`
	WebsiteURL  string `json:"website"`
}

// Compositions of the two structs above
type Coin struct {
	CoinData MintData
	IPFSData IPFS
}
