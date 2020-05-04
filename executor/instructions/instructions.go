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
	"SealEVM/environment"
	"SealEVM/evmErrors"
	"SealEVM/memory"
	"SealEVM/opcodes"
	"SealEVM/stack"
	"SealEVM/storageCache"
)

type instructionsSetting struct {
	stack   *stack.Stack
	memory  *memory.Memory
	storage storageCache.ICache
	context environment.Context
}

type opCodeAction func(setting *instructionsSetting) ([]byte, error)
type opCodeInstruction struct {
	doAction        opCodeAction
	minStackDepth   int
	enabled         bool
}

type IInstructions interface {
	Execute(code opcodes.OpCode) ([]byte, error)
}

var instructionTable [256]opCodeInstruction

func (i *instructionsSetting) Execute(code opcodes.OpCode) ([]byte, error) {
	instruction := instructionTable[code]
	if !instruction.enabled {
		return nil, evmErrors.InvalidOpCode(byte(code))
	}

	if !i.stack.CheckMinDepth(instruction.minStackDepth) {
		return nil, evmErrors.StackUnderFlow
	}

	return instructionTable[int(code)].doAction(i)
}

func Load()  {
	loadArithmetic()
	loadBitOperations()
	loadComparision()
}

func New(stack *stack.Stack, memory *memory.Memory, storage storageCache.ICache, context environment.Context) IInstructions {
	is := &instructionsSetting{
		stack:   stack,
		memory:  memory,
		storage: storage,
		context: context,
	}
	return is
}
