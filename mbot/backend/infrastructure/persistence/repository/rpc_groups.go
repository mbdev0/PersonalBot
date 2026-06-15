package repository

import (
	"context"
	"database/sql"
	"personal_bot/backend/infrastructure/persistence/mapper"
	"personal_bot/backend/infrastructure/persistence/models"
	rpcgroups "personal_bot/backend/internal/core/rpc_groups"
)

type RPCGroup struct {
	db *sql.DB
}

func NewRPCGroupRepository(db *sql.DB) *RPCGroup {
	return &RPCGroup{db: db}
}

func (rg *RPCGroup) GetDashboard(ctx context.Context) (rpcgroups.RPCGroupDashboard, error) {
	query := `SELECT id, group_name, json_array_length(rpc_groups) AS rpc_count, creation_time FROM rpc_groups`

	rows, err := rg.db.QueryContext(ctx, query)
	if err != nil {
		return rpcgroups.RPCGroupDashboard{}, err
	}

	var rpcDashboard rpcgroups.RPCGroupDashboard

	for rows.Next() {
		var rpcGroup rpcgroups.RPCGroupDashboardRow
		err := rows.Scan(&rpcGroup.Id, &rpcGroup.Name, &rpcGroup.Number, &rpcGroup.CreationTime)
		if err != nil {
			return rpcgroups.RPCGroupDashboard{}, err
		}

		rpcDashboard = append(rpcDashboard, rpcGroup)
	}
	return rpcDashboard, nil
}
func (rg *RPCGroup) Add(ctx context.Context, rpcGroup rpcgroups.RPCGroupPost) (rpcgroups.RPCGroup, error) {
	mappedRpcGroup, err := mapper.MapRpcGroupPostToRepository(rpcGroup)
	if err != nil {
		return rpcgroups.RPCGroup{}, err
	}

	query := "INSERT INTO rpc_groups (group_name, rpc_groups, creation_time) VALUES (?,?,?)"

	lastInsertId, err := execTx(ctx, rg.db, query, Fields{
		mappedRpcGroup.GroupName,
		mappedRpcGroup.RpcGroups,
		mappedRpcGroup.CreationTime,
	})
	if err != nil {
		return rpcgroups.RPCGroup{}, err
	}

	return rg.GetById(ctx, lastInsertId)
}

func (rg *RPCGroup) GetById(ctx context.Context, id int64) (rpcgroups.RPCGroup, error) {
	query := "SELECT * FROM rpc_groups where id=?"
	rows := rg.db.QueryRowContext(ctx, query, id)

	var rpcGroup models.RpcGroupRepository
	err := rows.Scan(&rpcGroup.Id, &rpcGroup.GroupName, &rpcGroup.RpcGroups, &rpcGroup.CreationTime)
	if err != nil {
		return rpcgroups.RPCGroup{}, err
	}

	mappedRpcGroup, err := mapper.MapRepositoryToRpcGroup(rpcGroup)
	if err != nil {
		return rpcgroups.RPCGroup{}, err
	}

	return mappedRpcGroup, nil
}

func (rg *RPCGroup) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM rpc_groups where id=?"

	_, err := execTx(ctx, rg.db, query, Fields{id})
	return err
}

func (rg *RPCGroup) Update(ctx context.Context, id int64, rcpGroupUpdate rpcgroups.RPCGroupPut) (rcpGroup rpcgroups.RPCGroup, err error) {
	query := `UPDATE rpc_groups SET group_name=?, rpc_groups=? WHERE id=?`

	mappedRpcGroupPut, err := mapper.MapRpcGroupPutToRepository(rcpGroupUpdate)
	if err != nil {
		return
	}

	_, err = execTx(ctx, rg.db, query, Fields{
		mappedRpcGroupPut.GroupName, mappedRpcGroupPut.RpcGroups, id,
	})
	if err != nil {
		return
	}

	return rg.GetById(ctx, id)
}
