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

func loadArithmetic() {
	instructionTable[opcodes.STOP] = opCodeInstruction {
		doAction:       stopAction,
		minStackDepth:  0,
		enabled:        true,
	}

	instructionTable[opcodes.ADD] = opCodeInstruction {
		doAction:       addAction,
		minStackDepth:  2,
		enabled:        true,
	}

	instructionTable[opcodes.MUL] = opCodeInstruction {
		doAction:       mulAction,
		minStackDepth:  2,
		enabled:        true,
	}

	instructionTable[opcodes.SUB] = opCodeInstruction {
		doAction:       subAction,
		minStackDepth:  2,
		enabled:        true,
	}

	instructionTable[opcodes.DIV] = opCodeInstruction {
		doAction:       divAction,
		minStackDepth:  2,
		enabled:        true,
	}

	instructionTable[opcodes.SDIV] = opCodeInstruction {
		doAction:       sDivAction,
		minStackDepth:  2,
		enabled:        true,
	}

	instructionTable[opcodes.MOD] = opCodeInstruction {
		doAction:       modAction,
		minStackDepth:  2,
		enabled:        true,
	}

	instructionTable[opcodes.SMOD] = opCodeInstruction {
		doAction:       sModAction,
		minStackDepth:  2,
		enabled:        true,
	}

	instructionTable[opcodes.ADDMOD] = opCodeInstruction {
		doAction:       addModAction,
		minStackDepth:  3,
		enabled:        true,
	}

	instructionTable[opcodes.MULMOD] = opCodeInstruction {
		doAction:       mulModAction,
		minStackDepth:  3,
		enabled:        true,
	}

	instructionTable[opcodes.EXP] = opCodeInstruction {
		doAction:       expAction,
		minStackDepth:  2,
		enabled:        true,
	}

	instructionTable[opcodes.SIGNEXTEND] = opCodeInstruction {
		doAction:       signExtendAction,
		minStackDepth:  2,
		enabled:        true,
	}
}

func stopAction(_ *instructionsContext) ([]byte, error) {
	return nil, nil
}

func addAction(ctx *instructionsContext) ([]byte, error) {
	x, _ := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.Add(x)
	return nil, nil
}

func mulAction(ctx *instructionsContext) ([]byte, error) {
	x, _ := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.Mul(x)
	return nil, nil
}

func subAction(ctx *instructionsContext) ([]byte, error) {
	x, _ := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.Set(x.Sub(y).Int)
	return nil, nil
}

func divAction(ctx *instructionsContext) ([]byte, error) {
	x, _ := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.Set(x.Div(y).Int)
	return nil, nil
}

func sDivAction(ctx *instructionsContext) ([]byte, error) {
	x, _ := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.Set(x.SDiv(y).Int)
	return nil, nil
}

func modAction(ctx *instructionsContext) ([]byte, error) {
	x, _ := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.Set(x.Mod(y).Int)
	return nil, nil
}

func sModAction(ctx *instructionsContext) ([]byte, error) {
	x, _ := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.Set(x.SMod(y).Int)
	return nil, nil
}

func addModAction(ctx *instructionsContext) ([]byte, error) {
	x, _ := ctx.stack.Pop()
	y, _ := ctx.stack.Pop()
	m := ctx.stack.Peek()

	m.Set(x.AddMod(y, m).Int)
	return nil, nil
}

func mulModAction(ctx *instructionsContext) ([]byte, error) {
	x, _ := ctx.stack.Pop()
	y, _ := ctx.stack.Pop()
	m := ctx.stack.Peek()

	m.Set(x.MulMod(y, m).Int)
	return nil, nil
}

func expAction(ctx *instructionsContext) ([]byte, error) {
	x, _ := ctx.stack.Pop()
	e := ctx.stack.Peek()

	e.Set(x.Exp(e).Int)
	return nil, nil
}

func signExtendAction(ctx *instructionsContext) ([]byte, error) {
	x, _ := ctx.stack.Pop()
	b := ctx.stack.Peek()

	b.Set(x.SignExtend(b).Int)
	return nil, nil
}
