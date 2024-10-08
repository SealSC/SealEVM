package gasSetting

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
)

func gasOfLog(topicCnt uint64) DynamicGasCalculator {
	return func(
		contract *environment.Contract,
		stx *stack.Stack,
		mem *memory.Memory,
		store *storage.Storage,
	) (uint64, uint64, error) {
		mOffset := stx.PeekPos(0)
		size := stx.PeekPos(1)

		expandSize, gasCost, err := mem.CalculateMallocSizeAndGas(mOffset, size)
		if err != nil {
			return 0, gasCost, err
		}
		return expandSize, 375 + 375*topicCnt + gasCost, nil
	}
}
