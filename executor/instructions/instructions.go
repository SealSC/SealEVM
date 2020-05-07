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
	"SealEVM/evmInt256"
	"SealEVM/memory"
	"SealEVM/stack"
	"SealEVM/storageCache"
)

type instructionsContext struct {
	stack       *stack.Stack
	memory      *memory.Memory
	storage     *storageCache.StorageCache
	environment environment.Context

	vm              interface{}
	pc              uint64
	lastReturn      []byte
	gasRemaining    *evmInt256.Int
	closureExec     ClosureExecute
}

type opCodeAction func(ctx *instructionsContext) ([]byte, error)
type opCodeInstruction struct {
	action            opCodeAction
	requireStackDepth int
	willIncreaseStack int
	enabled           bool
	jumps             bool
	returns           bool
	finished          bool
}

type IInstructions interface {
	ExecuteContract() ([]byte, error)
}

var instructionTable [256]opCodeInstruction

func (i *instructionsContext) ExecuteContract() ([]byte, error) {
	i.pc = 0
	contract := i.environment.Contract

	//todo: check if program is precompiled or nil contract
	var ret []byte
	var err error = nil

	for {
		opCode := contract.Code[i.pc]
		instruction := instructionTable[opCode]
		if !instruction.enabled {
			return nil, evmErrors.InvalidOpCode(opCode)
		}

		err := i.stack.CheckStackDepth(instruction.requireStackDepth, instruction.willIncreaseStack)
		if err != nil {
			return nil, err
		}

		ret, err = instruction.action(i)
		if err != nil {
			break
		}

		if !instruction.jumps {
			i.pc += 1
		}

		if instruction.finished {
			break
		}
	}

	return ret, err
}

func Load()  {
	loadStack()
	loadMemory()
	loadStorage()
	loadArithmetic()
	loadBitOperations()
	loadComparision()
	loadEnvironment()
	loadLog()
	loadMisc()
	loadClosure()
	loadPC()
}

func GetInstructionsTable() [256]opCodeInstruction {
	return instructionTable
}

func New(
	vm interface{},
	stack *stack.Stack,
	memory *memory.Memory,
	storage *storageCache.StorageCache,
	context environment.Context,
	closureExecute ClosureExecute) IInstructions {

	is := &instructionsContext{
		vm:             vm,
		stack:          stack,
		memory:         memory,
		storage:        storage,
		environment:    context,
		closureExec:    closureExecute,
	}
	return is
}
