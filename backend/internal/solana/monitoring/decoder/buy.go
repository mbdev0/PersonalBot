package decoder

import (
	"bytes"
	"encoding/binary"
	"errors"
)

func ExtractBuyAmountFromBuyInstruction(data []byte) (float64, error) {
	if len(data) < 8 {
		return 0, errors.New("data length is less than 8 bytes")
	}

	buf := bytes.NewBuffer(data[8:])

	amount, err := convertToUint64(buf)
	if err != nil {
		return 0, err
	}

	return (float64(amount) / 1e17) * 100, nil
}

func ExtractTokenAmountFromPfInstruction(data []byte) (uint64, error) {
	if len(data) < 8 {
		return 0, errors.New("data length is less than 8 bytes")
	}
	tokenAmountBytes := data[8:16]

	tokenAmount := binary.LittleEndian.Uint64(tokenAmountBytes)
	return tokenAmount, nil

}
