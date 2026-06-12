package repository

import (
	"context"
	"database/sql"
	"personal_bot/infrastructure/persistence/mapper"
	"personal_bot/infrastructure/persistence/models"
	"personal_bot/internal/core/position"
	"personal_bot/pkg/logger"
)

type PositionsRepository struct {
	db *sql.DB
}

func NewPositionRepository(db *sql.DB) *PositionsRepository {
	return &PositionsRepository{db: db}
}

func (pr *PositionsRepository) Add(ctx context.Context, position position.Position) error {
	query := `INSERT INTO positions (
					position_id, strategy_id, token_address,
					wallet_address, init_token_amount, tokens_remaining, remaining_cost_basis,
					finalized_profit, initial_amount, entry_price, mcap_entry, average_mcap_exit,
					total_marketcap_exit, number_of_sales, address_for_url
				)
				VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	mappedPosition, err := mapper.MapPositionToRepo(position)
	if err != nil {
		return err
	}

	_, err = execTx(ctx, pr.db, query, Fields{
		mappedPosition.PositionId, mappedPosition.StrategyId,
		mappedPosition.TokenAddress, mappedPosition.WalletAddress, mappedPosition.InitialTokenAmount,
		mappedPosition.TokensRemaining, mappedPosition.RemainingCostBasis, mappedPosition.FinalizedProfit,
		mappedPosition.InitialAmount, mappedPosition.EntryPrice, mappedPosition.McapEntry,
		mappedPosition.AverageMcapExit, mappedPosition.TotalMarketCapExit, mappedPosition.NumberOfSells,
		mappedPosition.AddressForUrl,
	})
	return err
}

func (pr *PositionsRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM positions WHERE position_id=?`
	_, err := execTx(ctx, pr.db, query, Fields{id})
	if err != nil {
		return err
	}

	return nil
}

func (pr *PositionsRepository) Update(ctx context.Context, position position.Position) error {
	query := `UPDATE positions SET
				strategy_id = ?,
				token_address = ?,
				wallet_address = ?,
				init_token_amount = ?,
				tokens_remaining = ?,
				remaining_cost_basis = ?,
				finalized_profit = ?,
				initial_amount = ?,
				entry_price = ?,
				mcap_entry = ?,
				average_mcap_exit = ?,
				total_marketcap_exit = ?,
				number_of_sales = ?,
				address_for_url = ?
				WHERE position_id = ?`

	mappedPosition, err := mapper.MapPositionToRepo(position)
	if err != nil {
		return err
	}

	_, err = execTx(ctx, pr.db, query, Fields{
		mappedPosition.StrategyId,
		mappedPosition.TokenAddress, mappedPosition.WalletAddress, mappedPosition.InitialTokenAmount,
		mappedPosition.TokensRemaining, mappedPosition.RemainingCostBasis, mappedPosition.FinalizedProfit,
		mappedPosition.InitialAmount, mappedPosition.EntryPrice, mappedPosition.McapEntry,
		mappedPosition.AverageMcapExit, mappedPosition.TotalMarketCapExit, mappedPosition.NumberOfSells,
		mappedPosition.AddressForUrl, mappedPosition.PositionId,
	})

	if err != nil {
		return err
	}

	return nil
}

func (pr *PositionsRepository) GetAll(ctx context.Context) ([]position.Position, error) {
	query := `
		SELECT
			position_id, strategy_id, token_address,
			wallet_address, init_token_amount, tokens_remaining, remaining_cost_basis,
			finalized_profit, initial_amount, entry_price, mcap_entry, average_mcap_exit,
			total_marketcap_exit, number_of_sales, address_for_url
		FROM positions
	`

	rows, err := pr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []position.Position
	for rows.Next() {
		var row models.PositionRepository
		err := rows.Scan(
			&row.PositionId, &row.StrategyId, &row.TokenAddress,
			&row.WalletAddress, &row.InitialTokenAmount, &row.TokensRemaining, &row.RemainingCostBasis,
			&row.FinalizedProfit, &row.InitialAmount, &row.EntryPrice, &row.McapEntry, &row.AverageMcapExit,
			&row.TotalMarketCapExit, &row.NumberOfSells, &row.AddressForUrl,
		)
		if err != nil {
			return nil, err
		}

		mappedPosition, err := mapper.MapRepoToPosition(row)
		if err != nil {
			return nil, err
		}

		positions = append(positions, mappedPosition)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return positions, nil
}

func (pr *PositionsRepository) handleRollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		logger.Error(err)
	}
}
