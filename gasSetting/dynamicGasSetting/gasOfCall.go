package dynamicGasSetting

import (
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/opcodes"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
	"github.com/SealSC/SealEVM/types"
)

type CallGas func(
	code opcodes.OpCode,
	availableGas uint64,
	stx *stack.Stack,
	mem *memory.Memory,
	store *storage.Storage,
) (uint64, uint64, uint64, error)

func gasSendWithCall(availableGas, baseGas, requestedGas uint64) uint64 {
	remainingGas := availableGas - baseGas
	allButOne64th := remainingGas - (remainingGas / 64)
	if requestedGas > allButOne64th {
		return allButOne64th
	}

	return requestedGas
}

func gasOfCall(
	code opcodes.OpCode,
	availableGas uint64,
	stx *stack.Stack,
	mem *memory.Memory,
	store *storage.Storage,
) (uint64, uint64, uint64, error) {
	var baseGas uint64
	var mOffset, size *evmInt256.Int

	if code == opcodes.CALL || code == opcodes.CALLCODE {
		mOffset = stx.PeekPos(5)
		size = stx.PeekPos(6)

		val := stx.PeekPos(2)
		addr := stx.PeekPos(1)

		if !val.IsZero() {
			baseGas += 9000
			if code == opcodes.CALL {
				if !store.ContractExist(types.Int256ToAddress(addr)) {
					baseGas += 25000
				}
			}
		}
	} else {
		mOffset = stx.PeekPos(4)
		size = stx.PeekPos(5)
	}

	baseGas += gasWithTouchedCheck(stx, 1, store.CachedContract)

	expSize, memCost, err := mem.CalculateMallocSizeAndGas(mOffset, size)
	if err != nil {
		return 0, baseGas, 0, err
	}

	baseGas += memCost

	reqGas := stx.PeekPos(0).Uint64()
	sendGas := gasSendWithCall(availableGas, baseGas, reqGas)

	return expSize, baseGas + sendGas, sendGas, nil
}
