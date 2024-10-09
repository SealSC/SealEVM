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

func loadArithmetic() {
	instructionTable[opcodes.STOP] = opCodeInstruction{
		action:   stopAction,
		enabled:  true,
		finished: true,
	}

	instructionTable[opcodes.ADD] = opCodeInstruction{
		action:            addAction,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.MUL] = opCodeInstruction{
		action:            mulAction,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.SUB] = opCodeInstruction{
		action:            subAction,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.DIV] = opCodeInstruction{
		action:            divAction,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.SDIV] = opCodeInstruction{
		action:            sDivAction,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.MOD] = opCodeInstruction{
		action:            modAction,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.SMOD] = opCodeInstruction{
		action:            sModAction,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.ADDMOD] = opCodeInstruction{
		action:            addModAction,
		requireStackDepth: 3,
		enabled:           true,
	}

	instructionTable[opcodes.MULMOD] = opCodeInstruction{
		action:            mulModAction,
		requireStackDepth: 3,
		enabled:           true,
	}

	instructionTable[opcodes.EXP] = opCodeInstruction{
		action:            expAction,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.SIGNEXTEND] = opCodeInstruction{
		action:            signExtendAction,
		requireStackDepth: 2,
		enabled:           true,
	}
}

func stopAction(_ *instructionsContext) ([]byte, error) {
	return nil, nil
}

func addAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.Add(x)
	return nil, nil
}

func mulAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.Mul(x)
	return nil, nil
}

func subAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.Set(x.Sub(y).Int)
	return nil, nil
}

func divAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.Set(x.Div(y).Int)
	return nil, nil
}

func sDivAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.Set(x.SDiv(y).Int)
	return nil, nil
}

func modAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.Set(x.Mod(y).Int)
	return nil, nil
}

func sModAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Pop()
	y := ctx.stack.Peek()

	y.Set(x.SMod(y).Int)
	return nil, nil
}

func addModAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Pop()
	y := ctx.stack.Pop()
	m := ctx.stack.Peek()

	m.Set(x.AddMod(y, m).Int)
	return nil, nil
}

func mulModAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Pop()
	y := ctx.stack.Pop()
	m := ctx.stack.Peek()

	m.Set(x.MulMod(y, m).Int)
	return nil, nil
}

func expAction(ctx *instructionsContext) ([]byte, error) {
	x := ctx.stack.Pop()
	e := ctx.stack.Peek()

	e.Set(x.Exp(e).Int)
	return nil, nil
}

func signExtendAction(ctx *instructionsContext) ([]byte, error) {
	b := ctx.stack.Pop()
	x := ctx.stack.Peek()

	x.Set(x.SignExtend(b).Int)
	return nil, nil
}
