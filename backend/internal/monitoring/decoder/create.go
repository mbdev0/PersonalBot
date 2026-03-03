package decoder

import (
	"bytes"
	"fmt"
	"personal_bot/internal/core/constants"
	"personal_bot/internal/monitoring/models"
	"personal_bot/pkg/logger"

	"github.com/near/borsh-go"
)

func decodeCreateInstruction(data []byte) (*models.DecodedCreateInstruction, error) {
	if len(data) < 8 {
		return nil, nil
	}

	isCreateV1 := bytes.Equal(data[:8], constants.CreateInstructionDiscriminator[:])
	isCreateV2 := bytes.Equal(data[:8], constants.CreateV2InstructionDiscriminator[:])
	isCreate := isCreateV1 || isCreateV2
	if !isCreate {
		return nil, nil
	}

	if isCreateV1 {
		decoded, err := parseDecodedInstructionV1(data[8:])
		logger.Information("in v1")
		return decoded, err
	} else if isCreateV2 {
		decoded, err := parseDecodedInstructionV2(data[8:])
		logger.Information("in v2")
		return decoded, err
	} else {
		return nil, fmt.Errorf("error whilst parsing data for create: %v", data)
	}

}

func parseDecodedInstructionV1(data []byte) (*models.DecodedCreateInstruction, error) {
	v1 := new(models.DecodedCreateInstructionV1)
	logger.Information(data)
	err := borsh.Deserialize(v1, data)
	logger.Error(err)
	if err != nil {
		return nil, err
	}

	return &models.DecodedCreateInstruction{
		Name:    v1.Name,
		Symbol:  v1.Symbol,
		IpfsUrl: v1.IpfsUrl,
	}, nil
}

func parseDecodedInstructionV2(data []byte) (*models.DecodedCreateInstruction, error) {
	v2 := new(models.DecodedCreateInstructionV2)
	logger.Information(data)
	err := borsh.Deserialize(v2, data)
	if err != nil {
		logger.Error("v2: ", err)
		return nil, err
	}

	return &models.DecodedCreateInstruction{
		Name:         v2.Name,
		Symbol:       v2.Symbol,
		IpfsUrl:      v2.IpfsUrl,
		Creator:      v2.Creator,
		IsMayhemMode: v2.IsMayhemMode,
	}, nil
}

func updateCoinFromDecodedInstruction(coin *models.Coin, instruction *models.DecodedCreateInstruction) {
	if instruction == nil {
		return
	}

	coin.CoinData = models.MintData{
		Name:    instruction.Name,
		Symbol:  instruction.Symbol,
		IpfsUrl: instruction.IpfsUrl,
	}
}
