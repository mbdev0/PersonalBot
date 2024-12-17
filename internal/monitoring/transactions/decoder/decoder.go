package decoder

type Decoder struct {
	decodingStragegy Decoding_Strategy
}

func (d *Decoder) SetDecodingStrategy(strategy Decoding_Strategy) {
	d.decodingStragegy = strategy
}

func (d *Decoder) Decode(data []byte) (interface{}, error) {
	return d.decodingStragegy.Decode(data)
}
