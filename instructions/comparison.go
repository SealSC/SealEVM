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
)

func loadComparision() {
	instructionTable[opcodes.LT] = opCodeInstruction{
		action:            ltAction,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.GT] = opCodeInstruction{
		action:            gtAction,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.SLT] = opCodeInstruction{
		action:            sltAction,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.SGT] = opCodeInstruction{
		action:            sgtAction,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.EQ] = opCodeInstruction{
		action:            eqAction,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.ISZERO] = opCodeInstruction{
		action:            isZeroAction,
		requireStackDepth: 1,
		enabled:           true,
	}
}

func ltAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Pop()
	y := ctx.stack.Peek()

	if x.LT(y) {
		y.SetUint64(1)
	} else {
		y.SetUint64(0)
	}
	return nil, nil
}

func gtAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Pop()
	y := ctx.stack.Peek()

	if x.GT(y) {
		y.SetUint64(1)
	} else {
		y.SetUint64(0)
	}
	return nil, nil
}

func sltAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Pop()
	y := ctx.stack.Peek()

	if x.SLT(y) {
		y.SetUint64(1)
	} else {
		y.SetUint64(0)
	}
	return nil, nil
}

func sgtAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Pop()
	y := ctx.stack.Peek()

	if x.SGT(y) {
		y.SetUint64(1)
	} else {
		y.SetUint64(0)
	}
	return nil, nil
}

func eqAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Pop()
	y := ctx.stack.Peek()

	if x.EQ(y) {
		y.SetUint64(1)
	} else {
		y.SetUint64(0)
	}
	return nil, nil
}

func isZeroAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Peek()

	if x.IsZero() {
		x.SetUint64(1)
	} else {
		x.SetUint64(0)
	}
	return nil, nil
}
