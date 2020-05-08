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
	"SealEVM/opcodes"
	"SealEVM/stack"
	"SealEVM/storageCache"
)

type DynamicGasCostSetting struct {
	EXPBytesCost  uint64
	SHA3ByteCost     uint64
	MemoryByteCost   uint64
	LogByteCost      uint64
}

type GasSetting struct {
	ActionConstCost [opcodes.MaxOpCodesCount] uint64
	NewAccountCost  uint64
	DynamicCost     DynamicGasCostSetting
}

func DefaultGasSetting() *GasSetting {
	gs := &GasSetting{}

	for i, _ := range gs.ActionConstCost {
		gs.ActionConstCost[i] = 2
	}

	gs.ActionConstCost[opcodes.EXP] = 10
	gs.ActionConstCost[opcodes.SHA3] = 30
	gs.ActionConstCost[opcodes.LOG0] = 375
	gs.ActionConstCost[opcodes.LOG1] = 375 * 2
	gs.ActionConstCost[opcodes.LOG2] = 375 * 3
	gs.ActionConstCost[opcodes.LOG3] = 375 * 4
	gs.ActionConstCost[opcodes.LOG4] = 375 * 5
	gs.ActionConstCost[opcodes.SLOAD] = 200
	gs.ActionConstCost[opcodes.SSTORE] = 5000
	gs.ActionConstCost[opcodes.SELFDESTRUCT] = 30000

	gs.ActionConstCost[opcodes.CREATE] = 32000
	gs.ActionConstCost[opcodes.CREATE2] = 32000

	gs.DynamicCost.EXPBytesCost = 50
	gs.DynamicCost.SHA3ByteCost = 6
	gs.DynamicCost.MemoryByteCost = 2
	gs.DynamicCost.LogByteCost = 8

	return gs
}

type ConstOpGasCostSetting [opcodes.MaxOpCodesCount] uint64

type instructionsContext struct {
	stack       *stack.Stack
	memory      *memory.Memory
	storage     *storageCache.StorageCache
	environment environment.Context

	vm              interface{}

	pc              uint64
	gasSetting      *GasSetting
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
	ExecuteContract() ([]byte, uint64, error)
	SetGasLimit(uint64)
}

var instructionTable [opcodes.MaxOpCodesCount]opCodeInstruction

//returns offset, size in type uint64
func (i *instructionsContext) memoryGasCostAndMalloc(offset *evmInt256.Int, size *evmInt256.Int) (uint64, uint64, uint64, error) {
	gasLeft := i.gasRemaining.Uint64()
	o, s, increased, err := i.memory.WillIncrease(offset, size)
	if err != nil {
		return o, s, gasLeft, err
	}

	gasCost := increased * i.gasSetting.DynamicCost.MemoryByteCost
	if gasLeft < gasCost {
		return o, s, gasLeft, evmErrors.OutOfGas
	}

	gasLeft -= gasCost
	i.gasRemaining.SetUint64(gasLeft)

	i.memory.Malloc(o, s)
	return o, s, gasLeft, err
}

func (i *instructionsContext) SetGasLimit(gasLimit uint64) {
	i.gasRemaining.SetUint64(gasLimit)
}

func (i *instructionsContext) ExecuteContract() ([]byte, uint64, error) {
	i.pc = 0
	contract := i.environment.Contract

	//todo: check if program is precompiled or nil contract
	var ret []byte
	var err error = nil

	for {
		opCode := contract.Code[i.pc]

		instruction := instructionTable[opCode]
		if !instruction.enabled {
			return nil, i.gasRemaining.Uint64(), evmErrors.InvalidOpCode(opCode)
		}

		err = i.stack.CheckStackDepth(instruction.requireStackDepth, instruction.willIncreaseStack)
		if err != nil {
			break
		}

		gasLeft := i.gasRemaining.Uint64()

		constCost := i.gasSetting.ActionConstCost[opCode]
		if gasLeft >= constCost {
			gasLeft -= constCost
			i.gasRemaining.SetUint64(gasLeft)
		} else {
			err = evmErrors.OutOfGas
			break
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

	return ret, i.gasRemaining.Uint64(), err
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

func GetInstructionsTable() [opcodes.MaxOpCodesCount]opCodeInstruction {
	return instructionTable
}

func New(
	vm interface{},
	stack *stack.Stack,
	memory *memory.Memory,
	storage *storageCache.StorageCache,
	context environment.Context,
	gasSetting *GasSetting,
	closureExecute ClosureExecute) IInstructions {

	is := &instructionsContext{
		vm:             vm,
		stack:          stack,
		memory:         memory,
		storage:        storage,
		environment:    context,
		closureExec:    closureExecute,
	}

	is.gasRemaining = evmInt256.FromBigInt(context.Transaction.GasLimit.Int)

	if gasSetting != nil {
		is.gasSetting = gasSetting
	} else {
		is.gasSetting = DefaultGasSetting()
	}

	return is
}
