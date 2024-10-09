package dynamicGasSetting

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
	"github.com/SealSC/SealEVM/types"
)

type touchedCheck func(address types.Address) bool

func gasWithTouchedCheck(stx *stack.Stack, addrPos uint, check touchedCheck) uint64 {
	addr := stx.PeekPos(addrPos)
	if check(types.Int256ToAddress(addr)) {
		return 100
	} else {
		return 2600
	}
}

func gasOfBalance(
	contract *environment.Contract,
	stx *stack.Stack,
	mem *memory.Memory,
	store *storage.Storage,
) (uint64, uint64, error) {
	return 0, gasWithTouchedCheck(stx, 0, store.CachedContract), nil
}

func gasOfExtCodeSize(
	contract *environment.Contract,
	stx *stack.Stack,
	mem *memory.Memory,
	store *storage.Storage,
) (uint64, uint64, error) {
	return 0, gasWithTouchedCheck(stx, 0, store.CachedContract), nil
}

func gasOfExtCodeHash(
	contract *environment.Contract,
	stx *stack.Stack,
	mem *memory.Memory,
	store *storage.Storage,
) (uint64, uint64, error) {
	return 0, gasWithTouchedCheck(stx, 0, store.CachedContract), nil
}
