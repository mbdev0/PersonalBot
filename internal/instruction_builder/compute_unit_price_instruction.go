package instructionbuilder

import computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"

func GetComputeUnitBudgetInstruction(computeUnitPrice uint64) *computebudget.Instruction {
	computeUnitPriceInstruction := computebudget.NewSetComputeUnitPriceInstruction(computeUnitPrice).Build()
	return computeUnitPriceInstruction
}
