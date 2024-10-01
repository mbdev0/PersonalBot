package transactions

import (
	"fmt"
	"pump_fun/internal/constants"
	"pump_fun/internal/logger"

	"github.com/gagliardetto/solana-go"
)

type InstructionDecoder func(accounts []*solana.AccountMeta, data []byte) (interface{}, error)

type CustomInstructionDecoderDecider struct {
	Decoders map[[8]byte]InstructionDecoder
}

func (c *CustomInstructionDecoderDecider) GetDecodeInstruction(discriminator [8]byte) (InstructionDecoder, error) {
	decoder, ok := c.Decoders[discriminator]
	if ok {
		return decoder, nil
	}
	return nil, fmt.Errorf("strategy not found")
}

func CustomInstructionDecoder(accounts []*solana.AccountMeta, data []byte) (interface{}, error) {
	decoders := CustomInstructionDecoderDecider{
		Decoders: map[[8]byte]InstructionDecoder{
			constants.CreateInstructionDiscriminator: CreateInstructionDecoder, //TODO: Add a buy instruction decoder in next PR
		},
	}

	// Convert []byte to [8]byte, TODO: Find a better fix for this
	var discriminator [8]byte
	copy(discriminator[:], data[:8])

	decoderStrategy, err := decoders.GetDecodeInstruction(discriminator)
	if err != nil {
		logger.Log(logger.LevelWarn, "Strategy not found for: "+string(discriminator[:]))
		return nil, err
	}
	return decoderStrategy(accounts, data)
}
