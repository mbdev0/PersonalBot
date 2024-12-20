package decoder

type Decoding_Strategy interface {
	Decode(data []byte) (interface{}, error)
}
