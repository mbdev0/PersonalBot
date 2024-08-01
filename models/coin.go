package models

type Coin struct {
}

type IPFS struct {
}

// Compositions of the two structs above
type TransactionData struct {
	CoinData Coin
	IPFSData IPFS
}
