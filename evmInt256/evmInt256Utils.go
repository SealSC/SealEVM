package evmInt256

import (
	"encoding/hex"
	"github.com/SealSC/SealEVM/common"
)

func EVMIntToHashBytes(i *Int) [common.HashLength]byte {
	iBytes := i.Bytes()
	iLen := len(iBytes)

	var hash [common.HashLength]byte
	if iLen > common.HashLength {
		copy(hash[:], iBytes[iLen-common.HashLength:])
	} else {
		copy(hash[common.HashLength-iLen:], iBytes)
	}

	return hash
}

func HashBytesToEVMInt(hash [common.HashLength]byte) (*Int, error) {

	i := New(0)
	i.SetBytes(hash[:])

	return i, nil
}

func HexToEVMInt(hStr string) *Int {
	bytes, err := hex.DecodeString(hStr)
	if err != nil {
		return nil
	}

	return BytesDataToEVMInt(bytes)
}

func BytesDataToEVMInt(data []byte) *Int {
	var bytes []byte
	srcLen := len(data)
	if srcLen > common.MaxIntBytes {
		bytes = data[:common.MaxIntBytes]
	} else {
		bytes = data
	}

	i := New(0)
	i.SetBytes(bytes)

	return i
}
