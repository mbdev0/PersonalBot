package repository

import (
	"context"
	"database/sql"
	"fmt"
	"personal_bot/infrastructure/persistence/mapper"
	"personal_bot/infrastructure/persistence/models"
	"personal_bot/internal/core/models/wallets"
)

type Wallet struct {
	db *sql.DB
}

func NewWalletRepository(conn *sql.DB) *Wallet {
	return &Wallet{db: conn}
}

func (w *Wallet) GetWallets(ctx context.Context) ([]wallets.SolanaWallet, error) {
	query := "SELECT * FROM crypto_wallets"
	rows, err := w.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	wallets := []wallets.SolanaWallet{}
	for rows.Next() {
		var wallet models.WalletRepository
		err := rows.Scan(&wallet.Id, &wallet.WalletName, &wallet.Chain, &wallet.PrivateKey)
		if err != nil {
			return nil, err
		}

		mappedWallet, err := mapper.WalletRepoToWallet(wallet)
		if err != nil {
			return nil, err
		}

		wallets = append(wallets, mappedWallet)
	}

	return wallets, nil
}

func (w *Wallet) GetWalletById(ctx context.Context, id string) (wallets.SolanaWallet, error) {
	query := "SELECT * FROM crypto_wallets where id = ?"
	row := w.db.QueryRowContext(ctx, query, id)

	var wallet models.WalletRepository
	err := row.Scan(&wallet.Id, &wallet.WalletName, &wallet.Chain, &wallet.PrivateKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return wallets.SolanaWallet{}, fmt.Errorf("wallet not found: %s", id)
		}
		return wallets.SolanaWallet{}, err
	}

	mappedRow, err := mapper.WalletRepoToWallet(wallet)
	if err != nil {
		return wallets.SolanaWallet{}, err
	}

	return mappedRow, nil
}

func (w *Wallet) GetWalletByName(ctx context.Context, name string) (wallets.SolanaWallet, error) {
	query := "SELECT * FROM crypto_wallets where wallet_name = ?"
	row := w.db.QueryRowContext(ctx, query, name)

	var wallet models.WalletRepository
	err := row.Scan(&wallet.Id, &wallet.WalletName, &wallet.Chain, &wallet.PrivateKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return wallets.SolanaWallet{}, fmt.Errorf("wallet not found: %s", name)
		}
		return wallets.SolanaWallet{}, err
	}

	mappedRow, err := mapper.WalletRepoToWallet(wallet)
	if err != nil {
		return wallets.SolanaWallet{}, err
	}

	return mappedRow, nil
}

func (w *Wallet) InsertWallets(ctx context.Context, walletRepo models.WalletRepository) (bool, error) {
	query := "INSERT INTO `crypto_wallets` VALUES (?,?,?,?)"

	insertResult, err := w.db.ExecContext(ctx, query, walletRepo.Id, walletRepo.WalletName, walletRepo.Chain, walletRepo.PrivateKey)
	if err != nil {
		return false, err
	}

	rowsAmountUpdated, err := insertResult.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAmountUpdated > 0, nil
}

func (w *Wallet) DeleteWallet(ctx context.Context, id string) (bool, error) {
	query := "DELETE FROM `crypto_wallets` where id = ?"
	tx, err := w.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})
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

func (w *Wallet) UpdateWallet(ctx context.Context, id string, wallet models.WalletRepository) (wallets.SolanaWallet, error) {
	query := "UPDATE crypto_wallets SET wallet_name = ?, chain = ?, private_key = ? WHERE id = ?"
	tx, err := w.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})
	if err != nil {
		return wallets.SolanaWallet{}, err
	}

	res, err := tx.ExecContext(ctx, query, wallet.WalletName, wallet.Chain, wallet.PrivateKey, id)
	if err != nil {
		tx.Rollback()
		return wallets.SolanaWallet{}, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return wallets.SolanaWallet{}, err
	}

	if rowsAffected > 1 {
		tx.Rollback()
		return wallets.SolanaWallet{}, fmt.Errorf("more than 1 rows were effected")
	}

	err = tx.Commit()
	if err != nil {
		return wallets.SolanaWallet{}, err
	}

	updatedWallet, err := w.GetWalletById(ctx, id)
	if err != nil {
		return wallets.SolanaWallet{}, fmt.Errorf("error whilst getting updated wallet: %v", err)
	}

	return updatedWallet, nil
}
