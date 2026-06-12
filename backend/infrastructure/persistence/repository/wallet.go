package repository

import (
	"context"
	"database/sql"
	"fmt"
	"personal_bot/infrastructure/persistence/mapper"
	"personal_bot/infrastructure/persistence/models"
	"personal_bot/internal/core/wallets"
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

	_, err := execTx(ctx, w.db, query, Fields{
		walletRepo.Id,
		walletRepo.WalletName,
		walletRepo.Chain,
		walletRepo.PrivateKey,
	})
	return err == nil, err
}

func (w *Wallet) DeleteWallet(ctx context.Context, id string) (bool, error) {
	query := "DELETE FROM `crypto_wallets` where id = ?"
	_, err := execTx(ctx, w.db, query, Fields{id})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (w *Wallet) UpdateWallet(ctx context.Context, id string, wallet models.WalletRepository) (wallets.SolanaWallet, error) {
	query := "UPDATE crypto_wallets SET wallet_name = ?, chain = ?, private_key = ? WHERE id = ?"
	_, err := execTx(ctx, w.db, query, Fields{wallet.WalletName, wallet.Chain, wallet.PrivateKey, id})

	if err != nil {
		return wallets.SolanaWallet{}, err
	}

	return w.GetWalletById(ctx, id)
}
