package rpcgroups

type RPCGroup struct {
	Id           int64
	Name         string
	Group        Group
	CreationTime string
}

type LoadedRPCGroup struct {
	RPCGroup   RPCGroup
	GroupIndex int
	References int
}

type RPCGroupPost struct {
	Name         string
	Group        Group
	CreationTime string
}

type RPCGroupPut struct {
	Name  string
	Group Group
}

type Group []GroupItem

type GroupItem struct {
	Http string `json:"http"`
	WS   string `json:"ws"`
}

type RPCGroupDashboard []RPCGroupDashboardRow

type RPCGroupDashboardRow struct {
	Id           int64
	Name         string
	Number       int64
	CreationTime string
}
