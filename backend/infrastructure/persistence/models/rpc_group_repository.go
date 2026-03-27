package models

type RpcGroupRepository struct {
	Id           int64  `db:"id"`
	GroupName    string `db:"group_name"`
	RpcGroups    string `db:"rpc_groups"`
	CreationTime string `db:"creation_time"`
}

type RpcGroupObject struct {
	Http string `json:"http"`
	Ws   string `json:"ws"`
}

type RpcGroupDashboard []RpcGroupDashboardRow

type RpcGroupDashboardRow struct {
	Name  string `db:"group_name"`
	Count int    `db:"count"`
}
