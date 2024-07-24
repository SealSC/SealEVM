package instructions

import (
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/opcodes"
)

func loadDencun() {
	instructionTable[opcodes.BLOBHASH] = opCodeInstruction{
		action:            blobHashAction,
		requireStackDepth: 1,
		enabled:           true,
	}
}

func blobHashAction(ctx *instructionsContext) ([]byte, error) {
	index := ctx.stack.Peek()
	hashLen := evmInt256.New(int64(len(ctx.environment.Transaction.BlobHashes)))
	if index.LT(hashLen) {
		blobHash := ctx.environment.Transaction.BlobHashes[index.Uint64()]
		index.SetBytes(blobHash[:])
	} else {
		index.SetUint64(0)
	}
	return nil, nil
}
