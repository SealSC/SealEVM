package types

import (
	"encoding/hex"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/utils"
)

const (
	SlotKeyBytesLen = 32
)

type SlotKey [HashBytesLen]byte

func Int256ToSlot(i *evmInt256.Int) SlotKey {
	var s SlotKey
	s.SetBytes(i.Bytes())
	return s
}

func (s SlotKey) Int256() *evmInt256.Int {
	return evmInt256.FromBytes(s[:])
}

func (s *SlotKey) SetBytes(b []byte) *SlotKey {
	utils.BytesCopy(s[:], b)
	return s
}

func (s SlotKey) String() string {
	return "0x" + hex.EncodeToString(s[:])
}

func (s SlotKey) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *SlotKey) UnmarshalText(val []byte) error {
	var data []byte
	err := utils.HexToBytes(val, &data, false)
	if err == nil {
		s.SetBytes(data)
	}

	return err
}
