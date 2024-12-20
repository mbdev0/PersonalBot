package decoder

import (
	"bytes"
	"pump_fun/internal/logger"
	"pump_fun/internal/models"
)

type CreateInstructionDecoder struct {
}

func (*CreateInstructionDecoder) Decode(data []byte) (models.DecodedInstruction, error) {
	var err error
	buf := bytes.NewBuffer(data[8:])

	var args models.DecodedInstruction

	args.Name, err = readStringWithLengthAtStart(buf)
	if err != nil {
		logger.Log(logger.LevelError, "Error reading the Name", logger.Error(err))
		return models.DecodedInstruction{}, err
	}

	args.Symbol, err = readStringWithLengthAtStart(buf)
	if err != nil {
		logger.Log(logger.LevelError, "Error reading the Symbol", logger.Error(err))
		return models.DecodedInstruction{}, err
	}

	args.IPFS_URL, err = readStringWithLengthAtStart(buf)
	if err != nil {
		logger.Log(logger.LevelError, "Error reading the URI", logger.Error(err))
		return models.DecodedInstruction{}, err
	}

	result := models.DecodedInstruction{
		Name:     args.Name,
		Symbol:   args.Symbol,
		IPFS_URL: args.IPFS_URL,
	}

	return result, nil
}
