package transactions

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gagliardetto/solana-go"
)

type CreateInstructionArgs struct {
    Name   string
    Symbol string
    Uri    string
}

func CustomInstructionDecoder(accounts []*solana.AccountMeta, data []byte) (interface{}, error) {
    var createInstructionID = [9]byte{24, 30, 200, 40, 5, 28, 7, 119, 24}

    if !bytes.Equal(data[:9], createInstructionID[:]) {
        return nil, errors.New("invalid instruction identifier")
    }

    buf := bytes.NewBuffer(data[9:])

    var args CreateInstructionArgs
    var err error

    args.Name, err = skipNullsReadStringAndTrimPadding(buf)
    if err != nil {
        return nil, err
    }

    args.Symbol, err = skipNullsReadStringAndTrimPadding(buf)
    if err != nil {
        return nil, err
    }

    args.Uri, err = skipNullsReadStringAndTrimPadding(buf)
    if err != nil {
        return nil, err
    }

    result := map[string]string{
        "Name":   args.Name,
        "Symbol": args.Symbol,
        "Uri":    args.Uri,
    }

    return result, nil
}

func skipNullsReadStringAndTrimPadding(buf *bytes.Buffer) (string, error) {
    if err := skipLeadingNullBytes(buf); err != nil {
        return "", err
    }

    str, err := readNullTerminatedString(buf)
    if err != nil {
        return "", err
    }

    return strings.TrimSpace(str), nil
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
