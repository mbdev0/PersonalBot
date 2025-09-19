package position

import (
	"context"
	"pump_fun/internal/core/position"
	"sync"
)

type PositionStreamer struct {
	//position_id => cancel
	activeStreams map[string]context.CancelFunc
	mu            *sync.Mutex
}

func NewPositionStreamer() *PositionStreamer {
	return &PositionStreamer{
		activeStreams: map[string]context.CancelFunc{},
		mu:            &sync.Mutex{},
	}
}

func (ps *PositionStreamer) OnPositionCreated(id string, pos *position.Position) {}

func (ps *PositionStreamer) OnPositionClosed(id string, pos *position.Position) {}

func (ps *PositionStreamer) OnPositionUpdated(id string, pos *position.Position) {}

func (ps *PositionStreamer) StopStreaming(id string) {

}
