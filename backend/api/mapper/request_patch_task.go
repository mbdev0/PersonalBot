package mapper

import (
	"fmt"
	"personal_bot/api/dto"
	"personal_bot/internal/core/models/wallets"
	"personal_bot/internal/core/tasks"
	"personal_bot/internal/solana/utils"

	"github.com/gagliardetto/solana-go"
)

func MapReqPatchToPatch(req dto.PatchRequestTask, typ string, wallet *wallets.SolanaWallet) (patch tasks.TaskPatch, err error) {
	switch typ {
	case string(dto.Buy):
		bp, err := mapReqPatchToBuyPatch(req, wallet)
		if err != nil {
			return nil, fmt.Errorf("error whilst mapping buy patch: %w", err)
		}
		return bp, nil
	case string(dto.Sell):
		sp, err := mapReqPatchToSellPatch(req, wallet)
		if err != nil {
			return nil, fmt.Errorf("error whilst mapping sell patch: %w", err)
		}
		return sp, nil
	default:
		return nil, fmt.Errorf("not recognised type of task: %s", req.Type)
	}
}

func mapReqPatchToBuyPatch(req dto.PatchRequestTask, wallet *wallets.SolanaWallet) (patch *tasks.BuyPatch, err error) {
	bp := &tasks.BuyPatch{}

	bp.Fee = req.BuyFee
	bp.ComputeUnit = req.ComputeUnits
	bp.Slippage = req.Slippage

	if req.BuyAmount != nil {
		amnt := req.BuyAmount
		buyAmount := utils.ConvertSolToLamport(*amnt)
		bp.Amount = buyAmount
	}

	if wallet != nil {
		bp.Wallet = wallet
	}

	if req.TokenAddress != nil {
		token, err := solana.PublicKeyFromBase58(*req.TokenAddress)
		if err != nil {
			return nil, solana.ErrInstructionDecoderNotFound
		}
		bp.Token = &token
	}

	return bp, nil
}

func mapReqPatchToSellPatch(req dto.PatchRequestTask, wallet *wallets.SolanaWallet) (patch *tasks.SellPatch, err error) {
	sp := &tasks.SellPatch{}

	sp.Fee = req.SellFee
	sp.ComputeUnit = req.ComputeUnits
	sp.Slippage = req.Slippage

	if req.SellAmount != nil {
		amnt := req.SellAmount
		sp.Amount = amnt

	}

	if wallet != nil {
		sp.Wallet = wallet
	}

	if req.TokenAddress != nil {
		token, err := solana.PublicKeyFromBase58(*req.TokenAddress)
		if err != nil {
			return nil, solana.ErrInstructionDecoderNotFound
		}
		sp.Token = &token
	}

	return sp, nil
}
