package constants

const (
	WebSocketReadLimit      = 65536
	LamportsConversion      = 1_000_000_000
	MicrolamportsToLamports = 1_000_000
	TokenAmountDecimals     = 1_000_000
	ShortPublicAddressInt   = 5
)

var (
	CreateInstructionDiscriminator   = [8]byte{24, 30, 200, 40, 5, 28, 7, 119}
	CreateV2InstructionDiscriminator = [8]byte{214, 144, 76, 236, 95, 139, 49, 180}
	BuyInstructionDiscriminator      = [8]byte{102, 6, 61, 18, 1, 218, 235, 234}
	SellInstructionDiscriminator     = [8]byte{51, 230, 133, 164, 1, 127, 131, 173}
	BondingCurveDiscriminator        = [8]byte{23, 183, 248, 55, 96, 216, 172, 96}
	TradeEventDiscrimnator           = [8]byte{189, 219, 127, 211, 78, 230, 97, 238}
)

const (
	Retries = 5
)
