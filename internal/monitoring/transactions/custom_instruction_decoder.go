package transactions

import (
	"errors"
	"pump_fun/internal/logger"

	"github.com/gagliardetto/solana-go"
)

type InstructionDecoder func(accounts []*solana.AccountMeta, data []byte) (interface{}, error)

type CustomInstructionDecoderDecider struct {
	Decoders map[[8]byte]InstructionDecoder
}

func (c *CustomInstructionDecoderDecider) GetDecodeInstruction(key [8]byte) (InstructionDecoder, error) {
	decoder, ok := c.Decoders[key]
	if ok {
		return decoder, nil
	}
	return nil, errors.New("strategy not found")
}

func CustomInstructionDecoder(accounts []*solana.AccountMeta, data []byte) (interface{}, error) {
	decoders := CustomInstructionDecoderDecider{
		Decoders: map[[8]byte]InstructionDecoder{
			{24, 30, 200, 40, 5, 28, 7, 119}: CreateInstructionDecoder, //TODO: Add a buy instruction decoder in next PR
		},
	}

	// Convert []byte to [8]byte, TODO: Find a better fix for this
	var key [8]byte
	copy(key[:], data[:8])

	decoderStrategy, err := decoders.GetDecodeInstruction(key)
	if err != nil {
		logger.Log(logger.LevelWarn, "Strategy not found for: "+string(key[:]), logger.Error(err))
		return nil, err
	}
	return decoderStrategy(accounts, data)
}
