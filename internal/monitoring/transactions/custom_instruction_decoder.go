package transactions

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gagliardetto/solana-go"
)

type CreateInstructionArgs struct {
    Name   string
    Symbol string
    Uri    string
}

func CustomInstructionDecoder(accounts []*solana.AccountMeta, data []byte) (interface{},error) {
    var createInstructionID = [8]byte{24, 30, 200, 40, 5, 28, 7, 119}
    var err error

    if !bytes.Equal(data[:8], createInstructionID[:]) {
        fmt.Println("invalid instruction identifier")
        return nil, err
    }

    buf := bytes.NewBuffer(data[8:])

    var args CreateInstructionArgs
    
    args.Name,err = readStringWithLengthAtStart(buf)
    if err != nil {
        return nil, err
    }

    args.Symbol,err = readStringWithLengthAtStart(buf)
    if err != nil {
        return nil, err
    }

    args.Uri,err = readStringWithLengthAtStart(buf)
    if err != nil {
        return nil, err
    }

    result := map[string]string{
        "Name":   args.Name,
        "Symbol": args.Symbol,
        "Uri":    args.Uri,
    }

    return result, err
}

func readStringWithLengthAtStart(buf *bytes.Buffer) (string, error) {
    lengthOfString, err := buf.ReadByte()
    if err != nil {
        return "", err
    }

    if err := skipLeadingNullBytes(buf); err != nil {
        return "", err
    }

    stringBytes := make([]byte, lengthOfString)
    _, err = buf.Read(stringBytes)
    if err != nil {
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

func readNullTerminatedString(buf *bytes.Buffer) (string, error) {
    str, err := buf.ReadString(0x00)
    if err != nil && err != io.EOF {
        return "", err
    }
    
    return str[:len(str)-1], nil
}
