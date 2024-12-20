package decoder

import "pump_fun/internal/models"

type Decoder struct {
	decodingStragegy Decoding_Strategy
}

func (d *Decoder) SetDecodingStrategy(strategy Decoding_Strategy) {
	d.decodingStragegy = strategy
}

func (d *Decoder) Decode(data []byte) (models.DecodedInstruction, error) {
	return d.decodingStragegy.Decode(data)
}
