package wallet

import (
	"context"
	"pump_fun/infrastructure/persistence/mapper"
	"pump_fun/infrastructure/persistence/repository"
	"pump_fun/internal/core/models/wallets"
)

type Service struct {
	walletRepo *repository.Wallet
}

func NewWalletService(walletRepo *repository.Wallet) *Service {
	return &Service{walletRepo: walletRepo}
}

func (s *Service) GetAll(ctx context.Context) ([]wallets.SolanaWallet, error) {
	wallets, err := s.walletRepo.GetWallets(ctx)
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func (s *Service) InsertWallet(ctx context.Context, wallet wallets.SolanaWallet) (bool, error) {
	mappedWallet := mapper.WalletToWalletRepo(wallet)
	isInserted, err := s.walletRepo.InsertWallets(ctx, mappedWallet)
	if err != nil {
		return false, err
	}

	return isInserted, nil
}
