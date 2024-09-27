package types

import (
	"encoding/hex"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/utils"
)

const (
	AddressBytesLen = 20
)

type Address [AddressBytesLen]byte

func (a Address) Int256() *evmInt256.Int {
	return evmInt256.FromBytes(a[:])
}

func (a *Address) SetBytes(b []byte) *Address {
	utils.BytesCopy(a[:], b)
	return a
}

func (a Address) String() string {
	return "0x" + hex.EncodeToString(a[:])
}

func (a Address) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *Address) UnmarshalText(val []byte) error {
	var data []byte
	err := utils.HexToBytes(val, &data, false)
	if err == nil {
		a.SetBytes(data)
	}

	return err
}
