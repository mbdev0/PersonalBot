package geyser

// TransactionNotification represents the root of the JSON structure.
type TransactionNotification struct {
	JsonRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

// Params represents the "params" field in the JSON.
type Params struct {
	Result       TransactionResult `json:"result"`
	Subscription int64             `json:"subscription"`
}

// TransactionResult represents the "result" field in the JSON.
type TransactionResult struct {
	Signature   string      `json:"signature"`
	Slot        int64       `json:"slot"`
	Transaction Transaction `json:"transaction"`
}

// Transaction represents the "transaction" field in the JSON.
type Transaction struct {
	Meta               Meta               `json:"meta"`
	TransactionDetails TransactionDetails `json:"transaction"`
	Version            interface{}        `json:"version"` // Can be "legacy" or an integer
}

// Meta represents the "meta" field in the JSON.
type Meta struct {
	ComputeUnitsConsumed int                `json:"computeUnitsConsumed"`
	Err                  interface{}        `json:"err"`
	Fee                  int64              `json:"fee"`
	InnerInstructions    []InnerInstruction `json:"innerInstructions"`
	LogMessages          []string           `json:"logMessages"`
	PostBalances         []int64            `json:"postBalances"`
	PostTokenBalances    []TokenBalance     `json:"postTokenBalances"`
	PreBalances          []int64            `json:"preBalances"`
	PreTokenBalances     []TokenBalance     `json:"preTokenBalances"`
	Rewards              interface{}        `json:"rewards"`
	Status               Status             `json:"status"`
}

// InnerInstruction represents an element in the "innerInstructions" array.
type InnerInstruction struct {
	Index        int           `json:"index"`
	Instructions []Instruction `json:"instructions"`
}

// Instruction represents an instruction in the "instructions" array.
type Instruction struct {
	Accounts    []string    `json:"accounts,omitempty"`
	Data        string      `json:"data,omitempty"`
	Parsed      interface{} `json:"parsed,omitempty"`
	Program     string      `json:"program,omitempty"`
	ProgramId   string      `json:"programId,omitempty"`
	StackHeight *int        `json:"stackHeight,omitempty"`
}

// Info represents the "info" field in a parsed instruction.
type Info struct {
	Amount         string   `json:"amount,omitempty"`
	Authority      string   `json:"authority,omitempty"`
	Destination    string   `json:"destination,omitempty"`
	Source         string   `json:"source,omitempty"`
	ExtensionTypes []string `json:"extensionTypes,omitempty"`
	Account        string   `json:"account,omitempty"`
	Mint           string   `json:"mint,omitempty"`
	Owner          string   `json:"owner,omitempty"`
	NewAccount     string   `json:"newAccount,omitempty"`
	Lamports       int64    `json:"lamports,omitempty"`
}

// TokenBalance represents a token balance in the "preTokenBalances" or "postTokenBalances" array.
type TokenBalance struct {
	AccountIndex  int           `json:"accountIndex"`
	Mint          string        `json:"mint"`
	Owner         string        `json:"owner"`
	ProgramId     string        `json:"programId"`
	UiTokenAmount UiTokenAmount `json:"uiTokenAmount"`
}

// UiTokenAmount represents the "uiTokenAmount" field in a token balance.
type UiTokenAmount struct {
	Amount         string  `json:"amount"`
	Decimals       int     `json:"decimals"`
	UiAmount       float64 `json:"uiAmount"`
	UiAmountString string  `json:"uiAmountString"`
}

// Status represents the "status" field in the JSON.
type Status struct {
	Ok interface{} `json:"Ok"`
}

// TransactionDetails represents the "transaction" field in the JSON.
type TransactionDetails struct {
	Message    Message  `json:"message"`
	Signatures []string `json:"signatures"`
}

// Message represents the "message" field in the transaction details.
type Message struct {
	AccountKeys         []AccountKey         `json:"accountKeys"`
	AddressTableLookups []AddressTableLookup `json:"addressTableLookups"`
	Instructions        []Instruction        `json:"instructions"`
	RecentBlockhash     string               `json:"recentBlockhash"`
}

// AccountKey represents an account key in the "accountKeys" array.
type AccountKey struct {
	Pubkey   string `json:"pubkey"`
	Signer   bool   `json:"signer"`
	Source   string `json:"source"`
	Writable bool   `json:"writable"`
}

// AddressTableLookup represents an address table lookup in the "addressTableLookups" array.
type AddressTableLookup struct {
	AccountKey      string `json:"accountKey"`
	ReadonlyIndexes []int  `json:"readonlyIndexes"`
	WritableIndexes []int  `json:"writableIndexes"`
}
