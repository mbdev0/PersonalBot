package repository

import (
	"context"
	"database/sql"
	"fmt"
	"personal_bot/infrastructure/persistence/mapper"
	"personal_bot/infrastructure/persistence/models"
	rpcgroups "personal_bot/internal/core/rpc_groups"
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

	tx, err := rg.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})

	if err != nil {
		tx.Rollback()
		return rpcgroups.RPCGroup{}, err
	}

	res, err := tx.ExecContext(ctx, query, mappedRpcGroup.GroupName, mappedRpcGroup.RpcGroups, mappedRpcGroup.CreationTime)
	if err != nil {
		tx.Rollback()
		return rpcgroups.RPCGroup{}, err
	}

	numAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return rpcgroups.RPCGroup{}, err
	}

	if numAffected != 1 {
		tx.Rollback()
		return rpcgroups.RPCGroup{}, fmt.Errorf("error whilst adding a new row, we found that more/less than 1 row got affected: %d", numAffected)
	}

	err = tx.Commit()
	if err != nil {
		return rpcgroups.RPCGroup{}, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return rpcgroups.RPCGroup{}, fmt.Errorf("error whilst getting last insert id: %d", lastInsertId)
	}

	latestAddedRpcGroup, err := rg.GetById(ctx, lastInsertId)
	if err != nil {
		return rpcgroups.RPCGroup{}, fmt.Errorf("error whilst finding the row for RPC group we just added error: %w", err)
	}

	return latestAddedRpcGroup, nil
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

	tx, err := rg.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})
	if err != nil {
		tx.Rollback()
		return err
	}

	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	numAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if numAffected > 1 {
		tx.Rollback()
		return fmt.Errorf("deleted too many rows, about to delete: %d rows", numAffected)
	}

	if numAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("wasn't able to find the id passed in: %d", id)
	}

	return tx.Commit()
}

func (rg *RPCGroup) Update(ctx context.Context, id int64, rcpGroupUpdate rpcgroups.RPCGroupPut) (rcpGroup rpcgroups.RPCGroup, err error) {
	query := `UPDATE rpc_groups SET group_name=?, rpc_groups=? WHERE id=?`

	mappedRpcGroupPut, err := mapper.MapRpcGroupPutToRepository(rcpGroupUpdate)
	if err != nil {
		return
	}

	tx, err := rg.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})
	if err != nil {
		return
	}

	res, err := tx.ExecContext(ctx, query, mappedRpcGroupPut.GroupName, mappedRpcGroupPut.RpcGroups, id)
	if err != nil {
		tx.Rollback()
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("we updated >1 || 0 rows - check the ID + database and try again")
		tx.Rollback()
		return
	}

	tx.Commit()

	return rg.GetById(ctx, id)
}
