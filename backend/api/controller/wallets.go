package controller

import (
	"context"
	"pump_fun/api/dto"
	"pump_fun/api/mapper"
	"pump_fun/internal/services/wallet"
	"pump_fun/pkg/logger"
)

type WalletsController struct {
	walletService *wallet.Service
}

func NewWalletController(walletService *wallet.Service) *WalletsController {
	return &WalletsController{walletService: walletService}
}

func (wc *WalletsController) GetWallets(ctx context.Context) ([]dto.ResponseWalletDto, error) {
	wallets, err := wc.walletService.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	logger.Information(wallets)

	res := []dto.ResponseWalletDto{}
	for _, wallet := range wallets {
		mappedWallet := mapper.MapWalletToDto(wallet)
		res = append(res, mappedWallet)
	}

	return res, nil
}

func (wc *WalletsController) InsertWallet(ctx context.Context, wallet dto.RequestWalletDto) (succeeded bool, err error) {
	mappedWallet, err := mapper.MapWalletDtoToWallet(wallet)
	if err != nil {
		return false, err
	}

	succeeded, err = wc.walletService.InsertWallet(ctx, mappedWallet)
	if err != nil {
		return false, err
	}

	return succeeded, nil

}

func (wc *WalletsController) GetWalletById() dto.ResponseWalletDto {
	return dto.ResponseWalletDto{}
}
