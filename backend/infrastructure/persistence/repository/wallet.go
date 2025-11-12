package repository

import (
	"context"
	"database/sql"
	"pump_fun/infrastructure/persistence/mapper"
	"pump_fun/infrastructure/persistence/models"
	"pump_fun/internal/core/models/wallets"
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

func (w *Wallet) GetWalletById(ctx context.Context, id string) (models.WalletRepository, error) {
	return models.WalletRepository{}, nil
}

func (w *Wallet) GetWalletByName(ctx context.Context, name string) (models.WalletRepository, error) {
	return models.WalletRepository{}, nil
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
