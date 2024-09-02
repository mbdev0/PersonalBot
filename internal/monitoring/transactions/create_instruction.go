package transactions

import (
	"bytes"
	"pump_fun/internal/logger"
	"pump_fun/internal/models"

	"github.com/gagliardetto/solana-go"
)

func CreateInstructionDecoder(accounts []*solana.AccountMeta, data []byte) (interface{}, error) {
	var err error

	buf := bytes.NewBuffer(data[8:])

	var args models.DecodedInstruction

	args.Name, err = readStringWithLengthAtStart(buf)
	if err != nil {
		logger.Log(logger.LevelError, "Error reading the Name", logger.Error(err))
		return nil, err
	}

	args.Symbol, err = readStringWithLengthAtStart(buf)
	if err != nil {
		logger.Log(logger.LevelError, "Error reading the Symbol", logger.Error(err))
		return nil, err
	}

	args.IPFS_URL, err = readStringWithLengthAtStart(buf)
	if err != nil {
		logger.Log(logger.LevelError, "Error reading the URI", logger.Error(err))
		return nil, err
	}

	result := map[string]string{
		"Name":     args.Name,
		"Symbol":   args.Symbol,
		"IPFS_URL": args.IPFS_URL,
	}

	return result, nil
}
