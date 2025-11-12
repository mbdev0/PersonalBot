package constants

const (
	WebSocketReadLimit      = 65536
	LamportsConversion      = 1_000_000_000
	MicrolamportsToLamports = 1_000_000
	TokenAmountDecimals     = 1_000_000
)

var (
	CreateInstructionDiscriminator = [8]byte{24, 30, 200, 40, 5, 28, 7, 119}
	BuyInstructionDiscriminator    = [8]byte{102, 6, 61, 18, 1, 218, 235, 234}
	SellInstructionDiscriminator   = [8]byte{51, 230, 133, 164, 1, 127, 131, 173}
)

const (
	Retries = 5
)
