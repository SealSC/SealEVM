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
	"SealEVM/evmInt256"
	"SealEVM/opcodes"
)

func loadPC() {
	instructionTable[opcodes.JUMP] = opCodeInstruction {
		action:         jumpAction,
		minStackDepth:  1,
		enabled:        true,
		jumps:          true,
	}

	instructionTable[opcodes.JUMPI] = opCodeInstruction {
		action:         jumpIAction,
		minStackDepth:  2,
		enabled:        true,
		jumps:          true,
	}

	instructionTable[opcodes.JUMPDEST] = opCodeInstruction {
		action:        jumpDestAction,
		minStackDepth: 0,
		enabled:       true,
	}

	instructionTable[opcodes.PC] = opCodeInstruction {
		action:        pcAction,
		minStackDepth: 1,
		enabled:       true,
	}
}

func jumpAction(ctx *instructionsContext) ([]byte, error) {
	target, _ := ctx.stack.Pop()
	nextPC := target.Uint64()

	validJump, err := ctx.environment.Contract.IsValidJump(nextPC)
	if validJump {
		ctx.pc = nextPC
	}

	return nil, err
}

func jumpIAction(ctx *instructionsContext) ([]byte, error) {
	target, _ := ctx.stack.Pop()
	condition, _ := ctx.stack.Pop()
	nextPC := target.Uint64()

	if condition.Sign() != 0 {
		validJump, err := ctx.environment.Contract.IsValidJump(nextPC)
		if validJump {
			ctx.pc = nextPC
		}
		return nil, err
	} else {
		ctx.pc += 1
		return nil, nil
	}
}

func jumpDestAction(_ *instructionsContext) ([]byte, error) {
	return nil, nil
}

func pcAction(ctx *instructionsContext) ([]byte, error) {
	i := evmInt256.New(0)
	i.SetUint64(ctx.pc)

	err := ctx.stack.Push(i)
	return nil, err
}
