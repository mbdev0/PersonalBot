package transaction_decoder

import (
	"bytes"
	"pump_fun/internal/constants"
	"pump_fun/internal/logger"
	"pump_fun/internal/models"
)

func DecodeCreateInstruction(data []byte) (*models.DecodedCreateInstruction, error) {
	isCreate := len(data) < 8 || !bytes.Equal(data[:8], constants.CreateInstructionDiscriminator[:])
	if isCreate {
		return nil, nil
	}

	buf := bytes.NewBuffer(data[8:])
	var args models.DecodedCreateInstruction
	var err error

	args.Name, err = readStringWithLengthAtStart(buf)
	if err != nil {
		logger.Error("Error reading the Name", err)
		return nil, err
	}

	args.Symbol, err = readStringWithLengthAtStart(buf)
	if err != nil {
		logger.Error("Error reading the Symbol", err)
		return nil, err
	}

	args.IPFS_URL, err = readStringWithLengthAtStart(buf)
	if err != nil {
		logger.Error("Error reading the IPFS_URL", err)
		return nil, err
	}

	return &args, nil
}

func UpdateCoinFromDecodedInstruction(coin *models.Coin, instruction *models.DecodedCreateInstruction) {
	if instruction == nil {
		return
	}

	coin.CoinData = models.MintData{
		Name:     instruction.Name,
		Symbol:   instruction.Symbol,
		IPFS_URL: instruction.IPFS_URL,
	}
}
