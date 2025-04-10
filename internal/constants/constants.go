package constants

const (
	PumpFunAPIEndPoint = "https://frontend-api-v3.pump.fun/"
	WebSocketReadLimit = 65536
)

var (
	CreateInstructionDiscriminator = [8]byte{24, 30, 200, 40, 5, 28, 7, 119}
	BuyInstructionDiscriminator    = [8]byte{102, 6, 61, 18, 1, 218, 235, 234}
)
