package response

type AccountSubscribeModel struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Subscription uint64 `json:"subscription"`
		Result       struct {
			Context struct {
				Slot int `json:"slot"`
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
