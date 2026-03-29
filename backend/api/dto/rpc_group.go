package dto

type RPCGroupDashboard []RPCGroupDashboardRow

type RPCGroupDashboardRow struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	GroupCount   int    `json:"count"`
	CreationTime string `json:"creation_time"`
}

type RPCGroup struct {
	Id           int64   `json:"id"`
	Name         string  `json:"name"`
	Group        []Group `json:"group"`
	CreationTime string  `json:"creation_time"`
}

type RPCGroupResponse struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Group        string `json:"group"`
	CreationTime string `json:"creation_time"`
}

type Group struct {
	Http string `json:"http"`
	Ws   string `json:"ws"`
}

// we will process the string internally in b.e.
type RPCGroupPush struct {
	Name  string `json:"name"`
	Group string `json:"group"`
}
