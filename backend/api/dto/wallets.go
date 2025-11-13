package dto

type RequestWalletDto struct {
	Wallet_name string `json:"wallet_name"`
	Chain       string `json:"chain"`
	Private_key string `json:"private_key"`
}

type ResponseWalletDto struct {
	Id         string `json:"id"`
	WalletName string `json:"wallet_name"`
	Chain      string `json:"chain"`
	PublicKey  string `json:"public_key"`
}

type RequestPatchWalletDto struct {
	Wallet_name *string `json:"wallet_name"`
	Chain       *string `json:"chain"`
	Private_key *string `json:"private_key"`
}
