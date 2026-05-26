package models

type AccountSubscribeModel struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Subscription uint64 `json:"subscription"`
		Result       struct {
			Context struct {
				Slot uint64 `json:"slot"`
			} `json:"context"`
			Value struct {
				Lamports   uint64   `json:"lamports"`
				Data       []string `json:"data"`
				Owner      string   `json:"owner"`
				Executable bool     `json:"executable"`
				RentEpoch  uint64   `json:"rentEpoch"`
				Space      int      `json:"space"`
			} `json:"value"`
		} `json:"result"`
	} `json:"params"`
}

type AccountNotification struct {
	JsonRPC string                    `json:"jsonrpc"`
	Method  string                    `json:"method"`
	Params  AccountNotificationParams `json:"params"`
}

type AccountNotificationParams struct {
	Subscription int64                     `json:"subscription"`
	Result       AccountNotificationResult `json:"result"`
}

type AccountNotificationResult struct {
	Context AccountNotificationContext `json:"context"`
	Value   AccountNotificationValue   `json:"value"`
}

type AccountNotificationContext struct {
	Slot uint64 `json:"slot"`
}

type AccountNotificationValue struct {
	Lamports   uint64                  `json:"lamports"`
	Data       AccountNotificationData `json:"data"`
	Owner      string                  `json:"owner"`
	Executable bool                    `json:"executable"`
	RentEpoch  uint64                  `json:"rentEpoch"`
	Space      int                     `json:"space"`
}

type AccountNotificationData struct {
	Program string                    `json:"program"`
	Parsed  AccountNotificationParsed `json:"parsed"`
	Space   int                       `json:"space"`
}

type AccountNotificationParsed struct {
	Info AccountNotificationInfo `json:"info"`
	Type string                  `json:"type"`
}

type AccountNotificationInfo struct {
	Extensions  []AccountNotificationExtension `json:"extensions"`
	IsNative    bool                           `json:"isNative"`
	Mint        string                         `json:"mint"`
	Owner       string                         `json:"owner"`
	State       string                         `json:"state"`
	TokenAmount AccountNotificationTokenAmount `json:"tokenAmount"`
}

type AccountNotificationExtension struct {
	Extension string `json:"extension"`
}

type AccountNotificationTokenAmount struct {
	Amount         string  `json:"amount"`
	Decimals       int     `json:"decimals"`
	UiAmount       float64 `json:"uiAmount"`
	UiAmountString string  `json:"uiAmountString"`
}
