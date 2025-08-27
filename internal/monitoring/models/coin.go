package models

type DecodedCreateInstruction struct {
	Name    string
	Symbol  string
	IpfsUrl string
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

// Coin Compositions of the two structs above
type Coin struct {
	CoinData MintData
	IPFSData IPFS
}
