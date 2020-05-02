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
		exec: stopAction,
	}

	iTable[opcodes.ADD] = opCodeInstruction {
		exec: addAction,
	}

	iTable[opcodes.MUL] = opCodeInstruction{
		exec: mulAction,
	}

	iTable[opcodes.SUB] = opCodeInstruction{
		exec: subAction,
	}

	iTable[opcodes.DIV] = opCodeInstruction{
		exec: divAction,
	}

	iTable[opcodes.SDIV] = opCodeInstruction{
		exec: sDivAction,
	}

	iTable[opcodes.MOD] = opCodeInstruction{
		exec: modAction,
	}

	iTable[opcodes.SMOD] = opCodeInstruction{
		exec: sModAction,
	}

	iTable[opcodes.ADDMOD] = opCodeInstruction{
		exec: addModAction,
	}

	iTable[opcodes.MULMOD] = opCodeInstruction{
		exec: mulModAction,
	}

	iTable[opcodes.EXP] = opCodeInstruction{
		exec: expAction,
	}

	iTable[opcodes.SIGNEXTEND] = opCodeInstruction{
		exec: signExtendAction,
	}
}

func stopAction(context interface{}) ([]byte, error) {
	return nil, nil
}

func addAction(context interface{}) ([]byte, error) {
	return nil, nil
}

func mulAction(context interface{}) ([]byte, error) {
	return nil, nil
}

func subAction(context interface{}) ([]byte, error) {
	return nil, nil
}

func divAction(context interface{}) ([]byte, error) {
	return nil, nil
}

func sDivAction(context interface{}) ([]byte, error) {
	return nil, nil
}

func modAction(context interface{}) ([]byte, error) {
	return nil, nil
}

func sModAction(context interface{}) ([]byte, error) {
	return nil, nil
}

func addModAction(context interface{}) ([]byte, error) {
	return nil, nil
}

func mulModAction(context interface{}) ([]byte, error) {
	return nil, nil
}

func expAction(context interface{}) ([]byte, error) {
	return nil, nil
}

func signExtendAction(context interface{}) ([]byte, error) {
	return nil, nil
}
