package models

type DecodedInstruction struct {
    Name        string
    Symbol 		string
    Uri         string
}

type MintData struct {
	Signature        string
	Name             string
	Ticker           string // Symbol
	IPFS_URL         string // URI
	TokenAddr        string
	CreatorAddr      string
	DevHoldingAmount float64
}

type IPFS struct {
	TelegramURL string
	TwitterURL  string
	WebsiteURL  string
	ImageURL    string
}

// Compositions of the two structs above
type Coin struct {
	CoinData MintData
	IPFSData IPFS
}
