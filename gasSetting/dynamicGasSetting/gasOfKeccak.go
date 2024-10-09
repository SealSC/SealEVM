package dynamicGasSetting

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
)

func gasOfKeccak(
	contract *environment.Contract,
	stx *stack.Stack,
	mem *memory.Memory,
	store *storage.Storage,
) (uint64, uint64, error) {
	var gasCost uint64 = 30

	offset := stx.PeekPos(0)
	dataSize := stx.PeekPos(1)

	gasCost += ((dataSize.Uint64() + 31) / 32) * 6

	expSize, memGasCost, err := mem.CalculateMallocSizeAndGas(offset, dataSize)
	if err != nil {
		return 0, 0, err
	}

	gasCost += memGasCost

	return expSize, gasCost, nil
}
