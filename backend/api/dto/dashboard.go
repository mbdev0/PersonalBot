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
	Data      TradingTaskResponse `json:"data"`     //the trading task itself holding all options
	Children  []ChildRow          `json:"children"` //a table row can/cannot have children depending on type
}

type ChildRow struct {
	Type      string `json:"type"`
	Id        int64  `json:"id"`
	WsMessage string `json:"wsMessage"`
	State     string `json:"state"`
	Data      ResponseTask
}

type DashboardResponseDto struct {
	Rows []TableRow `json:"rows"`
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
