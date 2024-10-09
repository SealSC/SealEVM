package dynamicGasSetting

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
)

func gasOfExp(
	contract *environment.Contract,
	stx *stack.Stack,
	mem *memory.Memory,
	store *storage.Storage,
) (uint64, uint64, error) {
	var gasCost uint64 = 10
	b := stx.PeekPos(1)

	if b.Sign() > 0 {
		gasCost += uint64((b.BitLen()+7)/8) * 50
	}

	return 0, gasCost, nil
}
