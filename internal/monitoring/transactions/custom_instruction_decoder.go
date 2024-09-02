package transactions

import (
	"bytes"
	"fmt"

	"pump_fun/internal/logger"

	"github.com/gagliardetto/solana-go"
)

// TODO: Remove this duplication
type CreateInstructionArgs struct {
	Name   string
	Symbol string
	Uri    string
}

func CustomInstructionDecoder(accounts []*solana.AccountMeta, data []byte) (interface{}, error) {
	var createInstructionID = [8]byte{24, 30, 200, 40, 5, 28, 7, 119}
	var err error

	if !bytes.Equal(data[:8], createInstructionID[:]) {
		logger.Log(logger.LevelInfo, "Invalid instruction identifier")
		return nil, fmt.Errorf("invalid instruction identifier")

	}

	buf := bytes.NewBuffer(data[8:])

	var args CreateInstructionArgs

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

	args.Uri, err = readStringWithLengthAtStart(buf)
	if err != nil {
		logger.Log(logger.LevelError, "Error reading the URI", logger.Error(err))
		return nil, err
	}

	result := map[string]string{
		"Name":   args.Name,
		"Symbol": args.Symbol,
		"Uri":    args.Uri,
	}

	return result, err
}
