package repository

import (
	"context"
	"database/sql"
	"fmt"
	"personal_bot/pkg/logger"
)

type Fields []any

func execTx(ctx context.Context, db *sql.DB, query string, fields Fields) (lastId int64, err error) {

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})

	defer func() {
		if err != nil {
			rollback(tx)
		}
	}()

	if err != nil {
		return 0, err
	}

	res, err := tx.ExecContext(ctx, query, fields...)
	if err != nil {
		return 0, err
	}

	numAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if numAffected != 1 {
		return 0, fmt.Errorf("error whilst adding a new row, we found that more/less than 1 row got affected: %d", numAffected)
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error whilst getting last insert id: %d", lastInsertId)
	}
	return lastInsertId, nil

}

func execTxBatch(ctx context.Context, db *sql.DB, query string, rows []Fields) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			rollback(tx)
		}
	}()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, fields := range rows {
		if _, err = stmt.ExecContext(ctx, fields...); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func rollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		logger.Error(err)
	}
}
