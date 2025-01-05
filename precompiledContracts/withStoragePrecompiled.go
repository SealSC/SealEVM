package precompiledContracts

import (
	"github.com/SealSC/SealEVM/evmErrors"
	"github.com/SealSC/SealEVM/types"
)

const (
	WithStoragePrecompiledStart uint64 = 0x20000
	WithStoragePrecompiledEnd   uint64 = 0x2FFFF
)

type IWithStoragePrecompiledContract interface {
	GasCost(addr types.Address, input []byte, dataBlock types.DataBlock) uint64
	Execute(addr types.Address, input []byte, dataBlock types.DataBlock) ([]byte, error)
}

var withStoragePrecompiledContracts = map[uint64]IWithStoragePrecompiledContract{}

func GetWithStoragePrecompiledContract(addr types.Address) IWithStoragePrecompiledContract {
	addrInt := addr.Int256()
	if !addrInt.IsUint64() {
		return nil
	}

	return withStoragePrecompiledContracts[addrInt.Uint64()]
}

func IsWithStoragePrecompiled(addr types.Address) bool {
	addrInt := addr.Int256()
	if !addrInt.IsUint64() {
		return false
	}

	return withStoragePrecompiledContracts[addrInt.Uint64()] != nil
}

func RegisterContractWithStorage(addr types.Address, c IWithStoragePrecompiledContract) error {
	addrInt := addr.Int256()
	if !addrInt.IsUint64() {
		return evmErrors.InvalidPrecompiledAddress(addr)
	}

	addrIdx := addrInt.Uint64()
	if addrIdx < WithStoragePrecompiledStart || addrIdx > WithStoragePrecompiledEnd {
		return evmErrors.InvalidPrecompiledAddress(addr)
	}

	if withStoragePrecompiledContracts[addrIdx] != nil {
		return evmErrors.DuplicatePrecompiledContract(addr)
	}

	withStoragePrecompiledContracts[addrIdx] = c
	return nil
}
