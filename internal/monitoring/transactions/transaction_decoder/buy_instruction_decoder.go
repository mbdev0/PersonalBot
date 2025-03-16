package transaction_decoder

import (
	"bytes"
	"pump_fun/internal/constants"
	"pump_fun/internal/models"
)

func DecodeBuyInstruction(coin *models.Coin, data []byte) error {
	if len(data) < 8 {
		return nil
	}

	if bytes.Equal(data[:8], constants.BuyInstructionDiscriminator[:]) {
		err := buy_decoder(data, coin)
		if err != nil {
			return err
		}
	}

	return nil
}

func buy_decoder(data []byte, coinModel *models.Coin) error {
	var err error
	buf := bytes.NewBuffer(data[8:])

	amount, err := convertToUint64(buf)
	if err != nil {
		return err
	}

	coinModel.CoinData.DevHoldingAmount = (float64(amount) / 100_000_000_000_000_000) * 100

	return nil
}
