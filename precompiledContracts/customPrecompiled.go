package precompiledContracts

import (
	"github.com/SealSC/SealEVM/evmErrors"
	"github.com/SealSC/SealEVM/types"
)

const (
	CustomPrecompiledStart uint64 = 0x10000
	CustomPrecompiledEnd   uint64 = 0x1FFFF
)

func RegisterContracts(addr types.Address, c PrecompiledContract) error {
	addrInt := addr.Int256()
	if !addrInt.IsInt64() {
		return evmErrors.InvalidPrecompiledAddress(addr)
	}

	addrIdx := addrInt.Uint64()
	if addrIdx < CustomPrecompiledStart || addrIdx > CustomPrecompiledEnd {
		return evmErrors.InvalidPrecompiledAddress(addr)
	}

	if contracts[addrIdx] != nil {
		return evmErrors.DuplicatePrecompiledContract(addr)
	}

	contracts[addrIdx] = c

	return nil
}
