package decoder

import (
	"bytes"
	"fmt"
	"pump_fun/pkg/logger"
)

func readStringWithLengthAtStart(buf *bytes.Buffer) (string, error) {
	lengthOfString, err := buf.ReadByte()
	if err != nil {
		logger.Error("Error reading length of string", err)
		return "", err
	}

	if err := skipLeadingNullBytes(buf); err != nil {
		logger.Error("Error skipping null bytes", err)
		return "", err
	}

	stringBytes := make([]byte, lengthOfString)
	_, err = buf.Read(stringBytes)
	if err != nil {
		logger.Error("Error reading string bytes", err)
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
			logger.Error("Error reading byte", err)
			return 0, err
		}
		num = num | uint64(b)<<uint(8*i)
	}
	return num, nil
}
