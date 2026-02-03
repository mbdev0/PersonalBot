package dto

type DashboardResponse struct {
	Strategies      []TradingTaskResponse    `json:"strategies"`
	TasksByStrategy map[int64][]ResponseTask `json:"tasksByStrategy"`
	ManualTasks     []ResponseTask           `json:"manualTasks"`
}

type TableRow struct {
	Type      string              `json:"type"`
	Id        int64               `json:"id"`
	WsMessage string              `json:"wsMessage"`
	State     string              `json:"state"`
	Data      TradingTaskResponse `json:"data"`
	Children  []ChildRow          `json:"children"`
}

type ChildRow struct {
	Type      string       `json:"type"`
	Id        int64        `json:"id"`
	WsMessage string       `json:"wsMessage"`
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
