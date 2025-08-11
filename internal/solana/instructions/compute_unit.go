package instructions

import (
	"pump_fun/internal/core/constants"
	"pump_fun/internal/solana/utils"
	"pump_fun/pkg/logger"

	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
)

func GetComputeUnitBudgetInstruction(buyFee float64, computeUnits uint32) *computebudget.Instruction {
	if computeUnits == 0 {
		logger.Error("computeUnits cannot be zero")
		return nil
	}
	totalLamports := utils.ConvertSolToLamport(buyFee)
	totalLamportsInt := totalLamports.Int64()

	pricePerUnitMicrolamports := (totalLamportsInt * constants.MicrolamportsToLamports) / int64(computeUnits)

	return computebudget.NewSetComputeUnitPriceInstruction(uint64(pricePerUnitMicrolamports)).Build()
}

func GetComputeUnitLimitInstruction(computeLimit uint32) *computebudget.Instruction {
	computeLimitInstruction := computebudget.NewSetComputeUnitLimitInstruction(computeLimit).Build()
	return computeLimitInstruction
}
