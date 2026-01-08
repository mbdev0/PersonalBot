package decoder

import (
	"bytes"
	"pump_fun/pkg/logger"
)

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
