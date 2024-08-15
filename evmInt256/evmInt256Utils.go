package evmInt256

import "github.com/SealSC/SealEVM/common"

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

func BytesDataToEVMIntHash(data []byte) *Int {
	var hashBytes []byte
	srcLen := len(data)
	if srcLen < common.HashLength {
		hashBytes = common.LeftPaddingSlice(data, common.HashLength)
	} else {
		hashBytes = data[:common.HashLength]
	}

	i := New(0)
	i.SetBytes(hashBytes)

	return i
}
