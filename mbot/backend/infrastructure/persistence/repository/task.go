package repository

import (
	"context"
	"database/sql"
	"fmt"
	"personal_bot/backend/infrastructure/persistence/mapper"
	"personal_bot/backend/infrastructure/persistence/models"
	"personal_bot/backend/internal/core/tasks"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (tr *TaskRepository) GetAll(ctx context.Context) ([]tasks.Task, error) {
	query := `
	SELECT t.id, t.task_type, t.slippage_percentage, t.compute_units, t.config, t.state, t.token, t.strategy_id, t.time_created, t.node_config,
		cw.id, cw.wallet_name, cw.chain, cw.private_key,
		p.id, p.program
		FROM tasks t
		INNER JOIN crypto_wallets cw ON cw.id = t.wallet_id
		INNER JOIN programs p ON p.id = t.program_id
	`
	rows, err := tr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	dbTasks := []tasks.Task{}
	for rows.Next() {
		task := models.TaskRow{}
		wallet := models.WalletRepository{}
		program := models.ProgramRepository{}
		err := rows.Scan(
			&task.Id, &task.TaskType, &task.Slippage,
			&task.ComputeUnits, &task.Config, &task.State,
			&task.Token, &task.StrategyId, &task.TimeCreatedUnix,
			&task.NodeConfig,
			&wallet.Id, &wallet.WalletName, &wallet.Chain, &wallet.PrivateKey,
			&program.Id, &program.Program,
		)

		if err != nil {
			return nil, err
		}

		mappedTask, err := mapper.MapRepoToTask(task, wallet, program)
		if err != nil {
			return nil, err
		}

		dbTasks = append(dbTasks, mappedTask)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dbTasks, nil
}

func (tr *TaskRepository) GetMaxId(ctx context.Context) int64 {
	query := `SELECT MAX(id) FROM tasks`
	row := tr.db.QueryRowContext(ctx, query)

	var id int64
	row.Scan(&id)

	return id
}

func (tr *TaskRepository) AddAllTasks(ctx context.Context, tasks []tasks.Task) (bool, error) {
	if len(tasks) == 0 {
		return true, nil
	}

	baseQuery := `
	INSERT INTO tasks VALUES (?,?,?,?,?,?,?,?,?,?,?, (SELECT id FROM programs WHERE program = ?))
	ON CONFLICT(id) DO UPDATE SET
	    task_type = excluded.task_type,
	    wallet_id = excluded.wallet_id,
	    slippage_percentage = excluded.slippage_percentage,
	    compute_units = excluded.compute_units,
	    config = excluded.config,
	    state = excluded.state,
	    token = excluded.token,
	    strategy_id = excluded.strategy_id,
	    time_created = excluded.time_created,
	    node_config = excluded.node_config,
	    program_id = excluded.program_id
	`

	rows := make([]Fields, len(tasks))
	for i, t := range tasks {
		mappedTask, err := mapper.MapTaskToRepo(t)
		if err != nil {
			return false, fmt.Errorf("failed to map task: %d", t.Id())
		}
		rows[i] = Fields{
			mappedTask.Id, mappedTask.TaskType, mappedTask.WalletId, mappedTask.Slippage,
			mappedTask.ComputeUnits, mappedTask.Config, mappedTask.StrategyId,
			mappedTask.State, mappedTask.Token, mappedTask.TimeCreatedUnix, mappedTask.NodeConfig,
			t.Program(),
		}
	}

	if err := execTxBatch(ctx, tr.db, baseQuery, rows); err != nil {
		return false, fmt.Errorf("error whilst adding tasks: %w", err)
	}

	return true, nil
}

func (tr *TaskRepository) DeleteAll(ctx context.Context) (_ bool, err error) {
	query := `DELETE from tasks`
	tx, err := tr.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			rollback(tx)
		}
	}()

	if _, err = tx.ExecContext(ctx, query); err != nil {
		return false, err
	}

	return true, tx.Commit()
}
