package dto

type DashboardResponse struct {
	Strategies      []TradingTaskResponse    `json:"strategies"`
	TasksByStrategy map[int64][]ResponseTask `json:"tasksByStrategy"`
	ManualTasks     []ResponseTask           `json:"manualTasks"`
}
