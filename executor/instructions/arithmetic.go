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

func loadArithmetic(iTable [256]opCodeInstruction) {
	iTable[opcodes.STOP] = opCodeInstruction{
		doAction: stopAction,
		minStackDepth: 0,
		enabled: true,
	}

	iTable[opcodes.ADD] = opCodeInstruction {
		doAction: addAction,
		minStackDepth: 2,
		enabled: true,
	}

	iTable[opcodes.MUL] = opCodeInstruction{
		doAction: mulAction,
		minStackDepth: 2,
		enabled: true,
	}

	iTable[opcodes.SUB] = opCodeInstruction{
		doAction: subAction,
		minStackDepth: 2,
		enabled: true,
	}

	iTable[opcodes.DIV] = opCodeInstruction{
		doAction: divAction,
		minStackDepth: 2,
		enabled: true,
	}

	iTable[opcodes.SDIV] = opCodeInstruction{
		doAction: sDivAction,
		minStackDepth: 2,
		enabled: true,
	}

	iTable[opcodes.MOD] = opCodeInstruction{
		doAction: modAction,
		minStackDepth: 2,
		enabled: true,
	}

	iTable[opcodes.SMOD] = opCodeInstruction{
		doAction: sModAction,
		minStackDepth: 2,
		enabled: true,
	}

	iTable[opcodes.ADDMOD] = opCodeInstruction{
		doAction: addModAction,
		minStackDepth: 3,
		enabled: true,
	}

	iTable[opcodes.MULMOD] = opCodeInstruction{
		doAction: mulModAction,
		minStackDepth: 3,
		enabled: true,
	}

	iTable[opcodes.EXP] = opCodeInstruction{
		doAction: expAction,
		minStackDepth: 2,
		enabled: true,
	}

	iTable[opcodes.SIGNEXTEND] = opCodeInstruction{
		doAction: signExtendAction,
		minStackDepth: 2,
		enabled: true,
	}
}

func stopAction(setting *instructionsSetting) ([]byte, error) {
	return nil, nil
}

func addAction(setting *instructionsSetting) ([]byte, error) {
	x, _ := setting.stack.Pop()
	y := setting.stack.Peek()

	y.Add(x)
	return nil, nil
}

func mulAction(setting *instructionsSetting) ([]byte, error) {
	x, _ := setting.stack.Pop()
	y := setting.stack.Peek()

	y.Mul(x)
	return nil, nil
}

func subAction(setting *instructionsSetting) ([]byte, error) {
	x, _ := setting.stack.Pop()
	y := setting.stack.Peek()

	y.Set(x.Sub(y).Int)
	return nil, nil
}

func divAction(setting *instructionsSetting) ([]byte, error) {
	x, _ := setting.stack.Pop()
	y := setting.stack.Peek()

	y.Set(x.Div(y).Int)
	return nil, nil
}

func sDivAction(setting *instructionsSetting) ([]byte, error) {
	x, _ := setting.stack.Pop()
	y := setting.stack.Peek()

	y.Set(x.SDiv(y).Int)
	return nil, nil
}

func modAction(setting *instructionsSetting) ([]byte, error) {
	x, _ := setting.stack.Pop()
	y := setting.stack.Peek()

	y.Set(x.Mod(y).Int)
	return nil, nil
}

func sModAction(setting *instructionsSetting) ([]byte, error) {
	x, _ := setting.stack.Pop()
	y := setting.stack.Peek()

	y.Set(x.SMod(y).Int)
	return nil, nil
}

func addModAction(setting *instructionsSetting) ([]byte, error) {
	x, _ := setting.stack.Pop()
	y, _ := setting.stack.Pop()
	m := setting.stack.Peek()

	m.Set(x.AddMod(y, m).Int)
	return nil, nil
}

func mulModAction(setting *instructionsSetting) ([]byte, error) {
	x, _ := setting.stack.Pop()
	y, _ := setting.stack.Pop()
	m := setting.stack.Peek()

	m.Set(x.MulMod(y, m).Int)
	return nil, nil
}

func expAction(setting *instructionsSetting) ([]byte, error) {
	x, _ := setting.stack.Pop()
	e := setting.stack.Peek()

	e.Set(x.Exp(e).Int)
	return nil, nil
}

func signExtendAction(setting *instructionsSetting) ([]byte, error) {
	x, _ := setting.stack.Pop()
	b := setting.stack.Peek()

	b.Set(x.SignExtend(b).Int)
	return nil, nil
}
