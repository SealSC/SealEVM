package instructions

import (
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/opcodes"
	"github.com/SealSC/SealEVM/storage"
)

func loadDencun() {
	instructionTable[opcodes.BLOBHASH] = opCodeInstruction{
		action:            blobHashAction,
		requireStackDepth: 1,
		enabled:           true,
	}

	instructionTable[opcodes.BLOBBASEFEE] = opCodeInstruction{
		action:            blobBaseFeeAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.TLOAD] = opCodeInstruction{
		action:            tLoadAction,
		requireStackDepth: 1,
		enabled:           true,
	}

	instructionTable[opcodes.TSTORE] = opCodeInstruction{
		action:            tStoreAction,
		requireStackDepth: 2,
		enabled:           true,
		isWriter:          true,
	}

	instructionTable[opcodes.MCOPY] = opCodeInstruction{
		action:            memCopyAction,
		requireStackDepth: 3,
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

func blobBaseFeeAction(ctx *instructionsContext) ([]byte, error) {
	bbf := ctx.environment.Block.BlobBaseFee
	if bbf == nil {
		bbf = evmInt256.New(0)
	} else {
		bbf = bbf.Clone()
	}

	ctx.stack.Push(bbf)
	return nil, nil
}

func tLoadAction(ctx *instructionsContext) ([]byte, error) {
	key := ctx.stack.Peek()
	val, err := ctx.storage.XLoad(ctx.environment.Contract.Namespace, key, storage.TStorage)
	if err != nil {
		return nil, err
	}
	key.Set(val.Int)
	return nil, nil
}

func tStoreAction(ctx *instructionsContext) ([]byte, error) {
	key := ctx.stack.Pop()
	val := ctx.stack.Pop()
	ctx.storage.XStore(ctx.environment.Contract.Namespace, key, val, storage.TStorage)
	return nil, nil
}

func memCopyAction(ctx *instructionsContext) ([]byte, error) {
	dst := ctx.stack.Pop()
	src := ctx.stack.Pop()
	length := ctx.stack.Pop()

	_, _, _, err := ctx.memoryGasCostAndMalloc(dst, length)
	if err != nil {
		return nil, err
	}

	err = ctx.memory.MCopy(dst.Uint64(), src.Uint64(), length.Uint64())
	return nil, err
}
