package transaction_decoder

import (
	"bytes"
	"fmt"
	"pump_fun/internal/logger"
)

func readStringWithLengthAtStart(buf *bytes.Buffer) (string, error) {
	lengthOfString, err := buf.ReadByte()
	if err != nil {
		logger.Log(logger.LevelError, "Error getting length of string", logger.Error(err))
		return "", err
	}

	if err := skipLeadingNullBytes(buf); err != nil {
		logger.Log(logger.LevelError, "Error skipping null bytes", logger.Error(err))
		return "", err
	}

	stringBytes := make([]byte, lengthOfString)
	_, err = buf.Read(stringBytes)
	if err != nil {
		logger.Log(logger.LevelError, "Error reading string", logger.Error(err))
		return "", err
	}

	str := string(stringBytes)
	return str, nil
}

func skipLeadingNullBytes(buf *bytes.Buffer) error {
	for {
		if buf.Len() == 0 {
			return fmt.Errorf("buffer is empty")
		}
		if buf.Bytes()[0] != '\x00' {
			break
		}
		buf.Next(1)
	}
	return nil
}

func convertToUint64(buf *bytes.Buffer) (uint64, error) {
	var num uint64
	for i := 0; i < 8; i++ {
		b, err := buf.ReadByte()
		if err != nil {
			logger.Log(logger.LevelError, "Error reading byte", logger.Error(err))
			return 0, err
		}
		num = num | uint64(b)<<uint(8*i)
	}
	return num, nil
}
