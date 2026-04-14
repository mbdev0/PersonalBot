package controller

import (
	"personal_bot/api/dto"
	"personal_bot/api/mapper"
	"personal_bot/internal/services/position"
	positionHub "personal_bot/internal/services/subscription_hub/position"
)

type PositionController struct {
	PositionService *position.Service
}

func NewPositionController(ps *position.Service) PositionController {
	return PositionController{
		PositionService: ps,
	}
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
func (pc *PositionController) GetBy(id int64) (dto.PositionDto, error) {
	position, err := pc.PositionService.GetById(id)
	if err != nil {
		return dto.PositionDto{}, err
	}

	mappedPosition := mapper.MapPositionToPositionDto(*position)
	return mappedPosition, nil
}

func (pc *PositionController) Subscribe(id int64, isInternalSub bool) (*positionHub.Subscription, error) {
	sub, err := pc.PositionService.Subscribe(id, isInternalSub)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (pc *PositionController) Unsubscribe(id int64, isInternalSub bool) error {
	err := pc.PositionService.Unsubscribe(id, isInternalSub)
	if err != nil {
		return err
	}

	return nil
}
