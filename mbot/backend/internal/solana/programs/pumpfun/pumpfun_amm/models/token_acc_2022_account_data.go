package models

type MintAccountInfo struct {
	Lamports   uint64          `json:"lamports"`
	Data       MintAccountData `json:"data"`
	Owner      string          `json:"owner"`
	Executable bool            `json:"executable"`
	RentEpoch  uint64          `json:"rentEpoch"`
	Space      int             `json:"space"`
}

type MintAccountData struct {
	Program string         `json:"program"`
	Parsed  MintParsedData `json:"parsed"`
	Space   int            `json:"space"`
}

type MintParsedData struct {
	Info MintInfo `json:"info"`
	Type string   `json:"type"`
}

type MintInfo struct {
	Decimals        int             `json:"decimals"`
	Extensions      []MintExtension `json:"extensions"`
	FreezeAuthority *string         `json:"freezeAuthority"`
	IsInitialized   bool            `json:"isInitialized"`
	MintAuthority   *string         `json:"mintAuthority"`
	Supply          string          `json:"supply"`
}

type MintExtension struct {
	Extension string             `json:"extension"`
	State     MintExtensionState `json:"state"`
}

type MintExtensionState struct {
	// metadataPointer fields
	Authority       *string `json:"authority"`
	MetadataAddress *string `json:"metadataAddress"`

	// tokenMetadata fields
	AdditionalMetadata []any   `json:"additionalMetadata"`
	Mint               *string `json:"mint"`
	Name               *string `json:"name"`
	Symbol             *string `json:"symbol"`
	UpdateAuthority    *string `json:"updateAuthority"`
	Uri                *string `json:"uri"`
}
