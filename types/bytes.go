package types

import (
	"encoding/hex"
	"encoding/json"
	"errors"
)

type Bytes []byte

// MarshalJSON implements the json.Marshaler interface.
func (b Bytes) MarshalJSON() ([]byte, error) {
	hexStr := "0x" + hex.EncodeToString(b)
	return json.Marshal(hexStr)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *Bytes) UnmarshalJSON(data []byte) error {
	var hexStr string
	if err := json.Unmarshal(data, &hexStr); err != nil {
		return err
	}

	if len(hexStr) < 2 || hexStr[:2] != "0x" {
		return errors.New("invalid hex string")
	}

	bytes, err := hex.DecodeString(hexStr[2:])
	if err != nil {
		return err
	}

	*b = bytes
	return nil
}
