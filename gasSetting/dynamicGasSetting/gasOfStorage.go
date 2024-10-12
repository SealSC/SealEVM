package dynamicGasSetting

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
	"github.com/SealSC/SealEVM/storage/cache"
	"github.com/SealSC/SealEVM/types"
)

type SStoreGas func(
	acc *environment.Account,
	stx *stack.Stack,
	store *storage.Storage,
) (gasCost uint64, err error)

func gasOfSLoad(
	contract *environment.Account,
	stx *stack.Stack,
	mem *memory.Memory,
	store *storage.Storage,
) (uint64, uint64, error) {
	slot := stx.PeekPos(0)
	org, _ := store.CachedData(contract.Address, types.Int256ToSlot(slot))
	if org == nil {
		return 0, 100, nil
	} else {
		return 0, 2100, nil
	}
}

func gasOfSStore(
	acc *environment.Account,
	stx *stack.Stack,
	mem *memory.Memory,
	store *storage.Storage,
) (uint64, uint64, error) {
	var gasCost uint64 = 0
	slot := stx.PeekPos(0)
	newVal := stx.PeekPos(1)
	org, current := store.CachedData(acc.Address, types.Int256ToSlot(slot))

	if org == nil {
		gasCost += 2100
		val, err := store.XLoad(acc.Address, types.Int256ToSlot(slot), cache.SStorage)
		if err != nil {
			return 0, gasCost, err
		}

		org = val
		current = val
	}

	if newVal.EQ(current) {
		gasCost += 100
		return 0, gasCost, nil
	}

	if current.EQ(org) {
		if org.IsZero() {
			gasCost += 20000
			return 0, gasCost, nil
		}

		gasCost += 2900
		return 0, gasCost, nil
	}

	gasCost += 100

	return 0, gasCost, nil
}
