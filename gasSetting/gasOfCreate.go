package gasSetting

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
	"github.com/SealSC/SealEVM/utils"
)

func gasOfCreate(isCreate2 bool) DynamicGasCalculator {
	return func(
		contract *environment.Contract,
		stx *stack.Stack,
		mem *memory.Memory,
		store *storage.Storage,
	) (uint64, uint64, error) {
		var gasCost uint64 = 32000
		var addrGenCost uint64 = 0

		mOffset := stx.PeekPos(1)
		size := stx.PeekPos(2)

		expSize, memCost, err := mem.CalculateMallocSizeAndGas(mOffset, size)
		if err != nil {
			return 0, gasCost, err
		}

		wordSize := utils.ToWordSize(size.Uint64())
		initCodeCost := 2 * wordSize
		if isCreate2 {
			addrGenCost = 6 * wordSize
		}

		return expSize, gasCost + memCost + initCodeCost + addrGenCost, nil
	}
}
