package evmInt256

import (
	"encoding/hex"
	"github.com/SealSC/SealEVM/utils"
)

func EVMIntToHashBytes(i *Int) [utils.HashLength]byte {
	iBytes := i.Bytes()
	iLen := len(iBytes)

	var hash [utils.HashLength]byte
	if iLen > utils.HashLength {
		copy(hash[:], iBytes[iLen-utils.HashLength:])
	} else {
		copy(hash[utils.HashLength-iLen:], iBytes)
	}

	return hash
}

func HashBytesToEVMInt(hash [utils.HashLength]byte) (*Int, error) {

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
	if srcLen > utils.MaxIntBytes {
		bytes = data[:utils.MaxIntBytes]
	} else {
		bytes = data
	}

	i := New(0)
	i.SetBytes(bytes)

	return i
}
