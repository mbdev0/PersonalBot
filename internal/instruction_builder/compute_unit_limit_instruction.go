package instructionbuilder

import computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"

func GetComputeUnitLimitInstruction(computeLimit uint32) *computebudget.Instruction {
	computeLimitInstruction := computebudget.NewSetComputeUnitLimitInstruction(computeLimit).Build()
	return computeLimitInstruction
}
