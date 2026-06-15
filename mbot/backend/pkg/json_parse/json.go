package jsonparse

import "encoding/json"

func Encode[T any](v T) ([]byte, error) {
	return json.Marshal(v)
}

func Decode[T any](data []byte) (T, error) {
	var v T
	err := json.Unmarshal(data, &v)
	return v, err
}
