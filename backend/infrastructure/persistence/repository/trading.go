package repository

import (
	"context"
	"database/sql"
	"fmt"
	"personal_bot/infrastructure/persistence/mapper"
	"personal_bot/infrastructure/persistence/models"
	"personal_bot/internal/core/strategies"
	"personal_bot/pkg/logger"
)

type TradingRepository struct {
	db *sql.DB
}

func NewTradingRepository(db *sql.DB) *TradingRepository {
	return &TradingRepository{db: db}
}

func (tr *TradingRepository) GetAllTasks(ctx context.Context) ([]strategies.Task, error) {
	query := `
	SELECT tt.id, tt.trading_type, tt.slippage, tt.compute_units, tt.config, tt.time_created,
	cw.id, cw.wallet_name, cw.chain, cw.private_key,
	rg.id, rg.group_name, rg.rpc_groups, rg.creation_time,
	p.id, p.program
		FROM trading_tasks tt
		INNER JOIN crypto_wallets cw ON cw.id = tt.wallet_id
		INNER JOIN rpc_groups rg ON rg.id = tt.rpc_group_id
		INNER JOIN programs p ON p.id = tt.program_id
	`
	rows, err := tr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	dbTasks := []strategies.Task{}
	for rows.Next() {
		task := models.TradingRow{}
		wallet := models.WalletRepository{}
		rpcGroup := models.RpcGroupRepository{}
		program := models.ProgramRepository{}
		err := rows.Scan(&task.Id, &task.TradingType, &task.Slippage, &task.ComputeUnits, &task.Config, &task.TimeCreatedUnix,
			&wallet.Id, &wallet.WalletName, &wallet.Chain, &wallet.PrivateKey,
			&rpcGroup.Id, &rpcGroup.GroupName, &rpcGroup.RpcGroups, &rpcGroup.CreationTime,
			&program.Id, &program.Program,
		)

		if err != nil {
			return nil, err
		}

		mappedTask, err := mapper.MapRepoToTradingTask(task, wallet, rpcGroup, program)
		if err != nil {
			return nil, err
		}

		dbTasks = append(dbTasks, mappedTask)
	}

	return dbTasks, nil
}

func (tr *TradingRepository) GetMaxId(ctx context.Context) int64 {
	query := `SELECT MAX(id) FROM trading_tasks`
	row := tr.db.QueryRowContext(ctx, query)

	var id int64
	row.Scan(&id)

	return id
}

func (tr *TradingRepository) AddAllTasks(ctx context.Context, tasks []strategies.Task) (bool, error) {
	if len(tasks) == 0 {
		return true, nil
	}

	baseQuery := `
		INSERT into trading_tasks values (?,?,?,?,?,?,?,?, (SELECT id FROM programs WHERE program = ?))
		ON CONFLICT(id) DO UPDATE SET
            trading_type = excluded.trading_type,
            wallet_id = excluded.wallet_id,
            slippage = excluded.slippage,
            compute_units = excluded.compute_units,
            config = excluded.config,
			time_created = excluded.time_created,
			rpc_group_id = excluded.rpc_group_id,
			program_id = excluded.program_id
	`

	tx, err := tr.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, baseQuery)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	for _, t := range tasks {
		mappedTask, err := mapper.MapTradingTaskToRepo(t)
		if err != nil {
			logger.Error(err)
			return false, fmt.Errorf("failed to map task: %d", t.StrategyTaskId())
		}

		_, err = stmt.ExecContext(ctx,
			mappedTask.Id, mappedTask.TradingType, mappedTask.WalletId, mappedTask.Slippage, mappedTask.ComputeUnits,
			mappedTask.Config, mappedTask.TimeCreatedUnix, mappedTask.RpcGroupId, t.GetProgram())
		if err != nil {
			logger.Error(err)
			return false, fmt.Errorf("error whilst executing data for task id: %d, error: %w", t.StrategyTaskId(), err)
		}
	}

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("error whilst commiting: %w", err)
	}

	return true, nil
}

func (tr *TradingRepository) DeleteAll(ctx context.Context) (bool, error) {
	query := `DELETE from trading_tasks`

	tx, err := tr.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (tr *TradingRepository) Delete(ctx context.Context, id int64) (bool, error) {
	query := "delete from trading_tasks where id = ?"
	tx, err := tr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})
	if err != nil {
		return false, err
	}

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}
