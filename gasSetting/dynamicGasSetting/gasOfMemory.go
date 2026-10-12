package dynamicGasSetting

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
)

func gasOfMemory(baseGas uint64, fixedSize *evmInt256.Int) CommonCalculator {
	return func(
		_ *environment.Account,
		stx *stack.Stack,
		mem *memory.Memory,
		_ *storage.Storage,
	) (uint64, uint64, error) {
		mOffset := stx.PeekPos(0)
		size := fixedSize
		if size == nil {
			size = stx.PeekPos(1)
		}

		expSize, gasCost, err := mem.CalculateMallocSizeAndGas(mOffset, size)
		return expSize, baseGas + gasCost, err
	}
}
