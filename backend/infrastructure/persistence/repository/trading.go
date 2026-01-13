package repository

import (
	"context"
	"database/sql"
	"pump_fun/internal/core/strategies"
)

type TradingRepository struct {
	db *sql.DB
}

func NewTradingRepository(db *sql.DB) *TradingRepository {
	return &TradingRepository{db: db}
}

func (tr *TradingRepository) AddAllTasks(tasks []strategies.Task, ctx context.Context) (bool, error) {
	return false, nil
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
