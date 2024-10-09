package dynamicGasSetting

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
	"github.com/SealSC/SealEVM/utils"
)

func gasOfCopyMem(mem *memory.Memory, offset *evmInt256.Int, size *evmInt256.Int, gasCost uint64) (uint64, uint64, error) {
	expSize, memGasCost, err := mem.CalculateMallocSizeAndGas(offset, size)

	if err != nil {
		return 0, 0, err
	}

	gasCost += utils.ToWordSize(size.Uint64()) * 3
	gasCost += memGasCost

	return expSize, gasCost, nil
}

func gasOfCopy(
	contract *environment.Contract,
	stx *stack.Stack,
	mem *memory.Memory,
	store *storage.Storage,
) (uint64, uint64, error) {
	var gasCost uint64 = 3

	var mOffset = stx.PeekPos(0)
	var size = stx.PeekPos(2)

	expSize, memGas, err := gasOfCopyMem(mem, mOffset, size, gasCost)
	if err != nil {
		return 0, gasCost, err
	}

	return expSize, gasCost + memGas, nil
}

func gasOfExtCodeCopy(
	contract *environment.Contract,
	stx *stack.Stack,
	mem *memory.Memory,
	store *storage.Storage,
) (uint64, uint64, error) {
	var gasCost = gasWithTouchedCheck(stx, 0, store.CachedContract)

	var mOffset = stx.PeekPos(1)
	var size = stx.PeekPos(3)

	expSize, memGas, err := gasOfCopyMem(mem, mOffset, size, gasCost)
	if err != nil {
		return 0, gasCost, nil
	}

	return expSize, gasCost + memGas, nil
}
