package solanerrors

// solanaInstructionErrors maps named Solana instruction errors to user-facing messages.
var solanaInstructionErrors = map[string]string{
	"InsufficientFundsForFee":     "Insufficient SOL for transaction fee",
	"InsufficientFundsForRent":    "Insufficient SOL for account rent",
	"AccountNotFound":             "Account not found",
	"InvalidAccountData":          "Invalid account data",
	"AccountDataTooSmall":         "Account data too small",
	"InvalidInstructionData":      "Invalid instruction data",
	"InvalidProgramId":            "Invalid program ID",
	"MissingRequiredSignature":    "Missing required signature",
	"UninitializedAccount":        "Account not initialized",
	"NotEnoughAccountKeys":        "Not enough account keys",
	"AccountAlreadyInitialized":   "Account already initialized",
	"IllegalOwner":                "Illegal account owner",
	"AccountBorrowFailed":         "Account borrow failed",
	"MaxSeedLengthExceeded":       "Max seed length exceeded",
	"InvalidSeeds":                "Invalid seeds",
	"BorshIoError":                "Serialization error",
	"AccountNotRentExempt":        "Account not rent exempt",
	"ProgramFailedToComplete":     "Program failed to complete",
	"ComputationalBudgetExceeded": "Compute budget exceeded",
	"BlockhashNotFound":           "Blockhash not found",
}
