package decoder

import "pump_fun/internal/models"

type Decoding_Strategy interface {
	Decode(data []byte) (models.DecodedInstruction, error)
}
