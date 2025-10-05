package mapper

import (
	"fmt"
	"pump_fun/api/dto"
	"pump_fun/internal/core/tasks"
	"pump_fun/internal/solana/utils"

	"github.com/gagliardetto/solana-go"
)

func MapReqPatchToPatch(req dto.PatchRequestTask, typ string) (patch tasks.TaskPatch, err error) {
	switch typ {
	case string(dto.Buy):
		bp, err := mapReqPatchToBuyPatch(req)
		if err != nil {
			return nil, fmt.Errorf("error whilst mapping buy patch: %w", err)
		}
		return bp, nil
	case string(dto.Sell):
		sp, err := mapReqPatchToSellPatch(req)
		if err != nil {
			return nil, fmt.Errorf("error whilst mapping sell patch: %w", err)
		}
		return sp, nil
	default:
		return nil, fmt.Errorf("not recognised type of task: %s", req.Type)
	}
}

func mapReqPatchToBuyPatch(req dto.PatchRequestTask) (patch *tasks.BuyPatch, err error) {
	bp := &tasks.BuyPatch{}

	bp.Fee = req.BuyFee
	bp.ComputeUnit = req.ComputeUnits
	bp.Slippage = req.Slippage

	if req.BuyAmount != nil {
		amnt := req.BuyAmount
		buyAmount := utils.ConvertSolToLamport(*amnt)
		bp.Amount = buyAmount
	}

	if req.WalletAddressPrivateKey != nil {
		wallet, err := solana.PrivateKeyFromBase58(*req.WalletAddressPrivateKey)
		if err != nil {
			return nil, err
		}
		bp.Wallet = &wallet
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

func mapReqPatchToSellPatch(req dto.PatchRequestTask) (patch *tasks.SellPatch, err error) {
	sp := &tasks.SellPatch{}

	sp.Fee = req.SellFee
	sp.ComputeUnit = req.ComputeUnits
	sp.Slippage = req.Slippage

	if req.SellAmount != nil {
		amnt := req.SellAmount
		sp.Amount = amnt

	}

	if req.WalletAddressPrivateKey != nil {
		wallet, err := solana.PrivateKeyFromBase58(*req.WalletAddressPrivateKey)
		if err != nil {
			return nil, err
		}
		sp.Wallet = &wallet
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
