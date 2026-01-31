package dto

type DashboardResponse struct {
	Strategies      []TradingTaskResponse    `json:"strategies"`
	TasksByStrategy map[int64][]ResponseTask `json:"tasksByStrategy"`
	ManualTasks     []ResponseTask           `json:"manualTasks"`
}

type TableRow struct {
	Type      string
	Id        int64
	WsMessage string
	State     string
	Data      TradingTaskResponse //the trading task itself holding all options
	Children  []ChildRow          //a table row can/cannot have children depending on type
}

type ChildRow struct {
	Type      string
	Id        int
	WsMessage string
	State     string
	Data      ResponseTask
}

type DashboardResponseDto struct {
	Rows []TableRow
}

//dashboard response should be

/*
{strategies: [
	{type: afk,
	id: 123,
	msg: null //->fill this in with w.s. messages
	state: string
	data: StrategyTask
	children:[
		{type: buyTask, id: 123, data: BuyTask} -> this needs to be the same format as the actual strategy -> ???}
		]
	},
	{type: buy,
	id: 123,
	data: BuyTask,
	msg: null //->fill this in with w.s. messages
	state: use data.state?

	},
	{type:sell,
	id: 123,
	data:SellTask,
	msg: null //->fill this in with w.s. messages
	state: use data.state?
	}
}
]}

*/
