package wallet

import (
	"context"
	"personal_bot/infrastructure/persistence/mapper"
	"personal_bot/infrastructure/persistence/repository"
	"personal_bot/internal/core/wallets"
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

func (s *Service) GetById(ctx context.Context, id string) (wallets.SolanaWallet, error) {
	wallet, err := s.walletRepo.GetWalletById(ctx, id)
	if err != nil {
		return wallets.SolanaWallet{}, err
	}

	return wallet, nil
}
func (s *Service) GetByName(ctx context.Context, name string) (wallets.SolanaWallet, error) {
	wallet, err := s.walletRepo.GetWalletByName(ctx, name)
	if err != nil {
		return wallets.SolanaWallet{}, err
	}

	return wallet, nil
}

func (s *Service) Delete(ctx context.Context, id string) (isDeleted bool, err error) {
	deleted, err := s.walletRepo.DeleteWallet(ctx, id)
	if err != nil {
		return deleted, err
	}

	return deleted, err
}

func (s *Service) Update(ctx context.Context, id string, wallet wallets.SolanaWallet) (wallets.SolanaWallet, error) {
	mappedWallet := mapper.WalletToWalletRepo(wallet)
	updatedWallet, err := s.walletRepo.UpdateWallet(ctx, id, mappedWallet)
	if err != nil {
		return wallets.SolanaWallet{}, err
	}

	return updatedWallet, nil
}
