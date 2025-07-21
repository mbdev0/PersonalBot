package instructionbuilder

import (
	"pump_fun/internal/constants"
	"pump_fun/internal/handlers"
	"pump_fun/pkg/logger"

	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
)

func GetComputeUnitBudgetInstruction(buyFee float64, computeUnits uint32) *computebudget.Instruction {
	if computeUnits == 0 {
		logger.Error("computeUnits cannot be zero")
		return nil
	}
	totalLamports := handlers.ConvertSolToLamport(buyFee)
	totalLamportsInt := totalLamports.Int64()

	pricePerUnitMicrolamports := (totalLamportsInt * constants.MicrolamportsToLamports) / int64(computeUnits)

	return computebudget.NewSetComputeUnitPriceInstruction(uint64(pricePerUnitMicrolamports)).Build()
}
