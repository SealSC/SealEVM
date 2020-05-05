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
	"SealEVM/opcodes"
)

func loadStorage() {
	instructionTable[opcodes.SLOAD] = opCodeInstruction{
		doAction: sLoadAction,
		minStackDepth: 1,
		enabled: true,
	}

	instructionTable[opcodes.SSTORE] = opCodeInstruction{
		doAction: sStoreAction,
		minStackDepth: 2,
		enabled: true,
	}
}

func sLoadAction(ctx *instructionsContext) ([]byte, error) {
	k := ctx.stack.Peek()

	v, err := ctx.storage.SLoad(ctx.environment.Contract.Address, k)
	if err != nil {
		return nil, err
	}

	k.Set(v.Int)
	return nil, nil
}

func sStoreAction(ctx *instructionsContext) ([]byte, error) {
	k, _ := ctx.stack.Pop()
	v, _ := ctx.stack.Pop()

	ctx.storage.SStore(ctx.environment.Contract.Address, k, v)
	return nil, nil
}
