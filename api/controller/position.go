package controller

import (
	"pump_fun/api/dto"
	"pump_fun/api/mapper"
	"pump_fun/internal/services/position"
	positionHub "pump_fun/internal/services/subscription_hub/position"
)

type PositionController struct {
	PositionService *position.Service
}

func (pc *PositionController) GetAll() []dto.PositionDto {
	allPositions := pc.PositionService.GetAll()

	mappedPositions := make([]dto.PositionDto, len(allPositions))
	for i, pos := range allPositions {
		mappedPos := mapper.MapPositionToPositionDto(pos)
		mappedPositions[i] = mappedPos
	}

	return mappedPositions
}

// no need to return a pointer since we're not copying/editing the values
func (pc *PositionController) GetBy(id string) (dto.PositionDto, error) {
	position, err := pc.PositionService.GetById(id)
	if err != nil {
		return dto.PositionDto{}, err
	}

	mappedPosition := mapper.MapPositionToPositionDto(*position)
	return mappedPosition, nil
}

func (pc *PositionController) Subscribe(id string) (*positionHub.Subscription, error) {
	sub, err := pc.PositionService.Subscribe(id)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (pc *PositionController) Unsubscribe(id string) error {
	err := pc.PositionService.Unsubscribe(id)
	if err != nil {
		return err
	}

	return nil
}
