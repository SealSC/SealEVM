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

func loadBitOperations() {
	instructionTable[opcodes.AND] = opCodeInstruction {
		action:        andAction,
		minStackDepth: 2,
		enabled:       true,
	}

	instructionTable[opcodes.OR] = opCodeInstruction {
		action:        orAction,
		minStackDepth: 2,
		enabled:       true,
	}

	instructionTable[opcodes.XOR] = opCodeInstruction {
		action:        xorAction,
		minStackDepth: 2,
		enabled:       true,
	}

	instructionTable[opcodes.NOT] = opCodeInstruction {
		action:        notAction,
		minStackDepth: 1,
		enabled:       true,
	}

	instructionTable[opcodes.BYTE] = opCodeInstruction {
		action:        byteAction,
		minStackDepth: 2,
		enabled:       true,
	}

	instructionTable[opcodes.SHL] = opCodeInstruction {
		action:        shlAction,
		minStackDepth: 2,
		enabled:       true,
	}

	instructionTable[opcodes.SHR] = opCodeInstruction {
		action:        shrAction,
		minStackDepth: 2,
		enabled:       true,
	}

	instructionTable[opcodes.SAR] = opCodeInstruction {
		action:        sarAction,
		minStackDepth: 2,
		enabled:       true,
	}

}

func andAction(ctx *instructionsContext) ([]byte, error) {
	x, _ := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.And(x)
	return nil, nil
}

func orAction(ctx *instructionsContext) ([]byte, error) {
	x, _ := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.Or(x)
	return nil, nil
}

func xorAction(ctx *instructionsContext) ([]byte, error) {
	x, _ := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.XOr(x)
	return nil, nil
}

func notAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Peek()

	x.Not(x)
	return nil, nil
}

func byteAction(ctx *instructionsContext) ([]byte, error) {
	i, _ := ctx.stack.Pop()
	x := ctx.stack.Peek()

	b := x.ByteAt(int(i.Uint64()))
	x.SetUint64(uint64(b))
	return nil, nil
}

func shlAction(ctx *instructionsContext) ([]byte, error) {
	s, _ := ctx.stack.Pop()
	x := ctx.stack.Peek()

	x.SHL(s.Uint64())
	return nil, nil
}

func shrAction(ctx *instructionsContext) ([]byte, error) {
	s, _ := ctx.stack.Pop()
	x := ctx.stack.Peek()

	x.SHR(s.Uint64())
	return nil, nil
}

func sarAction(ctx *instructionsContext) ([]byte, error) {
	s, _ := ctx.stack.Pop()
	x := ctx.stack.Peek()

	x.SAR(s.Uint64())
	return nil, nil
}
