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
	instructionTable[opcodes.AND] = opCodeInstruction{
		doAction: andAction,
		minStackDepth: 2,
		enabled: true,
	}

	instructionTable[opcodes.OR] = opCodeInstruction {
		doAction: orAction,
		minStackDepth: 2,
		enabled: true,
	}

	instructionTable[opcodes.XOR] = opCodeInstruction{
		doAction: xorAction,
		minStackDepth: 2,
		enabled: true,
	}

	instructionTable[opcodes.NOT] = opCodeInstruction{
		doAction: notAction,
		minStackDepth: 1,
		enabled: true,
	}

	instructionTable[opcodes.BYTE] = opCodeInstruction{
		doAction: byteAction,
		minStackDepth: 2,
		enabled: true,
	}

	instructionTable[opcodes.SHL] = opCodeInstruction{
		doAction: shlAction,
		minStackDepth: 2,
		enabled: true,
	}

	instructionTable[opcodes.SHR] = opCodeInstruction{
		doAction: shrAction,
		minStackDepth: 2,
		enabled: true,
	}

	instructionTable[opcodes.SAR] = opCodeInstruction{
		doAction: sarAction,
		minStackDepth: 2,
		enabled: true,
	}

}

func andAction(setting *instructionsSetting) ([]byte, error) {
	x, _ := setting.stack.Pop()
	y := setting.stack.Peek()

	y.And(x)
	return nil, nil
}

func orAction(setting *instructionsSetting) ([]byte, error) {
	x, _ := setting.stack.Pop()
	y := setting.stack.Peek()

	y.Or(x)
	return nil, nil
}

func xorAction(setting *instructionsSetting) ([]byte, error) {
	x, _ := setting.stack.Pop()
	y := setting.stack.Peek()

	y.XOr(x)
	return nil, nil
}

func notAction(setting *instructionsSetting) ([]byte, error) {
	x := setting.stack.Peek()

	x.Not(x)
	return nil, nil
}

func byteAction(setting *instructionsSetting) ([]byte, error) {
	i, _ := setting.stack.Pop()
	x := setting.stack.Peek()

	b := x.ByteAt(int(i.Uint64()))
	x.SetUint64(uint64(b))
	return nil, nil
}

func shlAction(setting *instructionsSetting) ([]byte, error) {
	s, _ := setting.stack.Pop()
	x := setting.stack.Peek()

	x.SHL(s.Uint64())
	return nil, nil
}

func shrAction(setting *instructionsSetting) ([]byte, error) {
	s, _ := setting.stack.Pop()
	x := setting.stack.Peek()

	x.SHR(s.Uint64())
	return nil, nil
}

func sarAction(setting *instructionsSetting) ([]byte, error) {
	s, _ := setting.stack.Pop()
	x := setting.stack.Peek()

	x.SAR(s.Uint64())
	return nil, nil
}
