package gasSetting

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
	"github.com/SealSC/SealEVM/types"
)

func gasOfSelfDestruct(
	contract *environment.Contract,
	stx *stack.Stack,
	mem *memory.Memory,
	store *storage.Storage,
) (uint64, uint64, error) {
	var gasCost uint64 = 5000

	receiver := types.Int256ToAddress(stx.PeekPos(0))

	balance, _ := store.Balance(contract.Address)
	if !balance.IsZero() && store.ContractEmpty(receiver) {
		gasCost += 25000
	}

	if !store.CachedContract(receiver) {
		gasCost += 2600
	}

	return 0, gasCost, nil
}
