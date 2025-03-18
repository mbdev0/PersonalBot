package transaction_decoder

import (
	"bytes"
	"errors"
	"pump_fun/internal/constants"
	"pump_fun/internal/models"
)

// TODO: Remove the below if my understanding is correct
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

// TODO: Remove the below if my understanding is correct
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

func ExtractBuyAmountFromBuyInstruction(data []byte) (float64, error) {
	if len(data) < 8 {
		return 0, errors.New("data length is less than 8 bytes")
	}

	buf := bytes.NewBuffer(data[8:])

	amount, err := convertToUint64(buf)
	if err != nil {
		return 0, err
	}

	return (float64(amount) / 100_000_000_000_000_000) * 100, nil
}
