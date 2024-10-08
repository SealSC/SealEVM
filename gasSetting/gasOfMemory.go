package gasSetting

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
)

func gasOfMemory(size *evmInt256.Int) DynamicGasCalculator {
	return func(
		contract *environment.Contract,
		stx *stack.Stack,
		mem *memory.Memory,
		store *storage.Storage,
	) (uint64, uint64, error) {
		mOffset := stx.PeekPos(0)
		if size == nil {
			size = stx.PeekPos(1)
		}
		return mem.CalculateMallocSizeAndGas(mOffset, size)
	}
}
