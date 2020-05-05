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
	"SealEVM/common"
	"SealEVM/evmInt256"
	"SealEVM/opcodes"
)

func loadMemory() {
	instructionTable[opcodes.MLOAD] = opCodeInstruction{
		doAction: mLoadAction,
		minStackDepth: 1,
		enabled: true,
	}

	instructionTable[opcodes.MSTORE] = opCodeInstruction{
		doAction: mStoreAction,
		minStackDepth: 2,
		enabled: true,
	}

	instructionTable[opcodes.MSTORE8] = opCodeInstruction{
		doAction: mStore8Action,
		minStackDepth: 2,
		enabled: true,
	}

	instructionTable[opcodes.MSIZE] = opCodeInstruction{
		doAction: mSizeAction,
		minStackDepth: 0,
		enabled: true,
	}

}

func mLoadAction(ctx *instructionsContext) ([]byte, error) {
	i := ctx.stack.Peek()
	offset := i.Uint64()

	bytes, err := ctx.memory.Map(offset, 32)
	if err != nil {
		return nil, err
	}

	i.SetBytes(bytes)
	return nil, nil
}

func mStoreAction(ctx *instructionsContext) ([]byte, error) {
	mOffset, _ := ctx.stack.Pop()
	v, _ := ctx.stack.Pop()

	valBytes := common.EVMIntToHashBytes(v)

	err := ctx.memory.Store(mOffset.Uint64(), valBytes[:])
	return nil, err
}

func mStore8Action(ctx *instructionsContext) ([]byte, error) {
	mOffset, _ := ctx.stack.Pop()
	v, _ := ctx.stack.Pop()
	valBytes := v.Uint64()

	err := ctx.memory.Set(mOffset.Uint64(), byte(valBytes & 0xff))
	return nil, err
}

func mSizeAction(ctx *instructionsContext) ([]byte, error) {
	err := ctx.stack.Push(evmInt256.New(ctx.memory.Size()))
	return nil, err
}
