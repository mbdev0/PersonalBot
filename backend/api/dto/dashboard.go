package dto

type TableRow struct {
	Program   string              `json:"program"`
	Type      string              `json:"type"`
	Id        int64               `json:"id"`
	WsMessage string              `json:"ws_message"`
	State     string              `json:"state"`
	Data      TradingTaskResponse `json:"data"`
	Children  []ChildRow          `json:"children"`
}

type ChildRow struct {
	Program   string       `json:"program"`
	Type      string       `json:"type"`
	Id        int64        `json:"id"`
	WsMessage string       `json:"ws_message"`
	State     string       `json:"state"`
	Data      ResponseTask `json:"data"`
	Children  []SellTasks  `json:"children"`
}

type SellTasks struct {
	Program   string       `json:"program"`
	Type      string       `json:"type"`
	Id        int64        `json:"id"`
	WsMessage string       `json:"ws_message"`
	State     string       `json:"state"`
	Data      ResponseTask `json:"data"`
}

type DashboardResponseDto struct {
	Rows []TableRow `json:"rows"`
}

func (dr *DashboardResponseDto) New() {
	dr.Rows = []TableRow{}
}

type RowType string

const (
	STRATEGY RowType = "strategy"
	TASK     RowType = "task"
)
