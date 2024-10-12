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
	"github.com/SealSC/SealEVM/opcodes"
	"github.com/SealSC/SealEVM/storage/cache"
	"github.com/SealSC/SealEVM/types"
)

func loadStorage() {
	instructionTable[opcodes.SLOAD] = opCodeInstruction{
		action:            sLoadAction,
		requireStackDepth: 1,
		enabled:           true,
	}

	instructionTable[opcodes.SSTORE] = opCodeInstruction{
		action:            sStoreAction,
		requireStackDepth: 2,
		enabled:           true,
		isWriter:          true,
	}
}

func sLoadAction(ctx *instructionsContext) ([]byte, error) {
	k := ctx.stack.Peek()

	slot := types.Int256ToSlot(k)
	v, err := ctx.storage.XLoad(ctx.environment.Address(), slot, cache.SStorage)
	if err != nil {
		return nil, err
	}

	k.Set(v.Int)
	return nil, nil
}

func sStoreAction(ctx *instructionsContext) ([]byte, error) {
	k := ctx.stack.Pop()
	v := ctx.stack.Pop()

	slot := types.Int256ToSlot(k)

	ctx.storage.XStore(ctx.environment.Address(), slot, v, cache.SStorage)
	return nil, nil
}
