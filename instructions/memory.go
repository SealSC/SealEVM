/*
 * Copyright 2020 The SealEVM Authors
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package instructions

import (
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/opcodes"
)

func loadMemory() {
	instructionTable[opcodes.MLOAD] = opCodeInstruction{
		action:            mLoadAction,
		requireStackDepth: 1,
		enabled:           true,
	}

	instructionTable[opcodes.MSTORE] = opCodeInstruction{
		action:            mStoreAction,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.MSTORE8] = opCodeInstruction{
		action:            mStore8Action,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.MSIZE] = opCodeInstruction{
		action:            mSizeAction,
		willIncreaseStack: 1,
		enabled:           true,
	}
}

func mLoadAction(ctx *instructionsContext) ([]byte, error) {
	mOffset := ctx.stack.Peek()

	//gas check
	offset, _, _, err := ctx.memoryGasCostAndMalloc(mOffset, evmInt256.New(32))
	if err != nil {
		return nil, err
	}

	bytes, err := ctx.memory.Map(offset, 32)
	if err != nil {
		return nil, err
	}

	mOffset.SetBytes(bytes)
	return nil, nil
}

func mStoreAction(ctx *instructionsContext) ([]byte, error) {
	mOffset := ctx.stack.Pop()
	v := ctx.stack.Pop()

	//gas check
	offset, _, _, err := ctx.memoryGasCostAndMalloc(mOffset, evmInt256.New(32))
	if err != nil {
		return nil, err
	}

	valBytes := evmInt256.EVMIntToHashBytes(v)
	err = ctx.memory.Store(offset, valBytes[:])
	return nil, err
}

func mStore8Action(ctx *instructionsContext) ([]byte, error) {
	mOffset := ctx.stack.Pop()
	v := ctx.stack.Pop()
	valBytes := v.Uint64()

	//gas check
	offset, _, _, err := ctx.memoryGasCostAndMalloc(mOffset, evmInt256.New(1))
	if err != nil {
		return nil, err
	}

	err = ctx.memory.Set(offset, byte(valBytes&0xff))
	return nil, err
}

func mSizeAction(ctx *instructionsContext) ([]byte, error) {
	ctx.stack.Push(evmInt256.New(ctx.memory.Size()))
	return nil, nil
}
