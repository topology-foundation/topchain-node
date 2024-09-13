package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// PinRequestKeyPrefix is the prefix to retrieve all PinRequest
	PinRequestKeyPrefix = "PinRequest/value/"
)

// PinRequestKey returns the store key to retrieve a PinRequest from the index fields
func PinRequestKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
