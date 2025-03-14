package transaction_decoder

import (
	"bytes"
	"pump_fun/internal/constants"
	"pump_fun/internal/logger"
	"pump_fun/internal/models"
)

func DecodeCreateInstruction(coin *models.Coin, data []byte) error {

	if len(data) < 8 {
		return nil
	}

	if bytes.Equal(data[:8], constants.CreateInstructionDiscriminator[:]) {
		err := create_decoder(data, coin)
		if err != nil {
			logger.Log(logger.LevelError, "Error decoding create instruction", logger.Error(err))
			return err
		}
	}

	return nil
}

func create_decoder(data []byte, coinModel *models.Coin) error {
	var err error
	buf := bytes.NewBuffer(data[8:])

	var args models.DecodedCreateInstruction

	args.Name, err = readStringWithLengthAtStart(buf)
	if err != nil {
		logger.Log(logger.LevelError, "Error reading the Name", logger.Error(err))
		return err
	}

	args.Symbol, err = readStringWithLengthAtStart(buf)
	if err != nil {
		logger.Log(logger.LevelError, "Error reading the Symbol", logger.Error(err))
		return err
	}

	args.IPFS_URL, err = readStringWithLengthAtStart(buf)
	if err != nil {
		logger.Log(logger.LevelError, "Error reading the URI", logger.Error(err))
		return err
	}

	coinModel.CoinData = models.MintData{
		Name:     args.Name,
		Symbol:   args.Symbol,
		IPFS_URL: args.IPFS_URL,
	}

	return nil
}
