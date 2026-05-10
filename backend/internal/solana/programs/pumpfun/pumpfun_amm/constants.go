package pumpfunamm

const (
	GlobalConfig            = "ADyA8hdefvWN2dbGGWFotbzWxrAvLW83WG6QCVXvJKqw"
	ProtocolFeeRecipient    = "7hTckgnGnLQR6sdH7YkqFTAA7VwTfYFaZ6EhEsU3saCX"
	EventAuthority          = "GS4CU59F31iL7aR2Q8zVS8DRrcRnXX1yjQ66TqNVQnaR"
	GlobalVolumeAccumulator = "C2aFPdENg4A2HQsmrd5rTw5TaYBX5Ku887cWjbFKtZpw"
	FeeConfig               = "5PHirr8joyTMp9JMm6nW7hNDVyEYdkzDqazxPD7RaTjx"
	BuyBackVault            = "5cjcW9wExnJJiqgLjq7DEG75Pm6JBgE1hNv4B2vHXUW6"
	BuyBackVaultWsol        = "GYH1Gae1wJytMSvMvw8JVcv7nuAbxi8i9erNVbERnzXd"
)

var (
	BuyExactQuoteInDiscriminator = []byte{198, 46, 21, 82, 180, 217, 232, 112}
	SellInstructionDiscriminator = []byte{51, 230, 133, 164, 1, 127, 131, 173}
	BuyEvent                     = []byte{103, 244, 82, 31, 44, 245, 119, 119}
	SellEvent                    = []byte{62, 47, 55, 10, 165, 3, 220, 42}
)
