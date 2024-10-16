package types

import (
	"encoding/hex"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/utils"
)

const (
	HashBytesLen = 32
)

type Hash [HashBytesLen]byte

func Int256ToHash(i *evmInt256.Int) Hash {
	var h Hash
	h.SetBytes(i.Bytes())
	return h
}

func (h Hash) Int256() *evmInt256.Int {
	return evmInt256.FromBytes(h[:])
}

func (h *Hash) SetBytes(b []byte) *Hash {
	*h = Hash{}
	utils.BytesCopy(h[:], b)
	return h
}

func (h Hash) String() string {
	return "0x" + hex.EncodeToString(h[:])
}

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

func (h *Hash) UnmarshalText(val []byte) error {
	var data []byte
	err := utils.HexToBytes(val, &data, false)
	if err == nil {
		h.SetBytes(data)
	}

	return err
}
