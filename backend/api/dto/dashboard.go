package dto

type DashboardResponse struct {
	Strategies      []TradingTaskResponse
	TasksByStrategy map[int64][]ResponseTask
	ManualTasks     []ResponseTask
}
