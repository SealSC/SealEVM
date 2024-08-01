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
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmErrors"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/opcodes"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
)

type DynamicGasCostSetting struct {
	EXPBytesCost   uint64
	SHA3ByteCost   uint64
	MemoryByteCost uint64
	LogByteCost    uint64
}

type GasSetting struct {
	ActionConstCost [opcodes.MaxOpCodesCount]uint64
	NewAccountCost  uint64
	DynamicCost     DynamicGasCostSetting
}

func DefaultGasSetting() *GasSetting {
	gs := &GasSetting{}

	for i := range gs.ActionConstCost {
		gs.ActionConstCost[i] = 3
	}

	gs.ActionConstCost[opcodes.EXP] = 10
	gs.ActionConstCost[opcodes.SHA3] = 30
	gs.ActionConstCost[opcodes.LOG0] = 375
	gs.ActionConstCost[opcodes.LOG1] = 375 * 2
	gs.ActionConstCost[opcodes.LOG2] = 375 * 3
	gs.ActionConstCost[opcodes.LOG3] = 375 * 4
	gs.ActionConstCost[opcodes.LOG4] = 375 * 5
	gs.ActionConstCost[opcodes.SLOAD] = 800
	gs.ActionConstCost[opcodes.SSTORE] = 5000
	gs.ActionConstCost[opcodes.SELFDESTRUCT] = 30000

	gs.ActionConstCost[opcodes.CREATE] = 32000
	gs.ActionConstCost[opcodes.CREATE2] = 32000

	gs.ActionConstCost[opcodes.TLOAD] = 100
	gs.ActionConstCost[opcodes.TSTORE] = 100

	gs.DynamicCost.EXPBytesCost = 50
	gs.DynamicCost.SHA3ByteCost = 6
	gs.DynamicCost.MemoryByteCost = 2
	gs.DynamicCost.LogByteCost = 8

	return gs
}

type ConstOpGasCostSetting [opcodes.MaxOpCodesCount]uint64

type instructionsContext struct {
	stack       *stack.Stack
	memory      *memory.Memory
	storage     *storage.Storage
	environment *environment.Context

	vm interface{}

	pc           uint64
	readOnly     bool
	gasSetting   *GasSetting
	lastReturn   []byte
	gasRemaining *evmInt256.Int
	closureExec  ClosureExecute
	exitOpCode   opcodes.OpCode
}

type opCodeAction func(ctx *instructionsContext) ([]byte, error)
type opCodeInstruction struct {
	action            opCodeAction
	requireStackDepth int
	willIncreaseStack int

	//flags
	enabled  bool
	jumps    bool
	isWriter bool
	returns  bool
	finished bool
}

type IInstructions interface {
	ExecuteContract() ([]byte, uint64, error)
	SetGasLimit(uint64)
	GetGasLeft() uint64
	GetGasSetting() *GasSetting
	SetReadOnly()
	IsReadOnly() bool
	ExitOpCode() opcodes.OpCode
}

var instructionTable [opcodes.MaxOpCodesCount]opCodeInstruction

// returns offset, size in type uint64
func (i *instructionsContext) memoryGasCostAndMalloc(offset *evmInt256.Int, size *evmInt256.Int) (uint64, uint64, uint64, error) {
	gasLeft := i.gasRemaining.Uint64()
	o, s, increased, err := i.memory.WillIncrease(*offset, *size)
	if err != nil {
		return o, s, gasLeft, err
	}

	if increased == 0 {
		return o, s, increased, nil
	}

	gasCost := increased * i.gasSetting.DynamicCost.MemoryByteCost
	if gasLeft < gasCost {
		return o, s, gasLeft, evmErrors.OutOfGas
	}

	gasLeft -= gasCost
	i.gasRemaining.SetUint64(gasLeft)

	i.memory.Malloc(increased)
	return o, s, gasLeft, err
}

func (i *instructionsContext) SetGasLimit(gasLimit uint64) {
	i.gasRemaining.SetUint64(gasLimit)
}

func (i *instructionsContext) SetReadOnly() {
	i.readOnly = true
}

func (i *instructionsContext) IsReadOnly() bool {
	return i.readOnly
}

func (i *instructionsContext) GetGasLeft() uint64 {
	return i.gasRemaining.Uint64()
}

func (i *instructionsContext) GetGasSetting() *GasSetting {
	return i.gasSetting
}

func (i *instructionsContext) ExitOpCode() opcodes.OpCode {
	return i.exitOpCode
}

func (i *instructionsContext) ExecuteContract() (ret []byte, gasRemaining uint64, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = evmErrors.Panicked(e.(error))
			gasRemaining = i.gasRemaining.Uint64()
		}
	}()

	i.pc = 0
	contract := i.environment.Contract

	if len(contract.Code) == 0 {
		return nil, i.gasRemaining.Uint64(), nil
	}

	for {
		opCode := contract.GetOpCode(i.pc)

		instruction := instructionTable[opCode]
		if !instruction.enabled {
			return nil, i.gasRemaining.Uint64(), evmErrors.InvalidOpCode(opCode)
		}

		if instruction.isWriter && i.readOnly {
			return nil, i.gasRemaining.Uint64(), evmErrors.WriteProtection
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

		if instruction.returns {
			i.lastReturn = ret
		}

		if err != nil {
			break
		}

		if !instruction.jumps {
			i.pc += 1
		}

		if instruction.finished {
			i.exitOpCode = opcodes.OpCode(opCode)
			break
		}
	}

	return ret, i.gasRemaining.Uint64(), err
}

func Load() {
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
	loadDencun()
}

func GetInstructionsTable() [opcodes.MaxOpCodesCount]opCodeInstruction {
	return instructionTable
}

func New(
	vm interface{},
	stack *stack.Stack,
	memory *memory.Memory,
	storage *storage.Storage,
	context *environment.Context,
	gasSetting *GasSetting,
	closureExecute ClosureExecute) IInstructions {

	is := &instructionsContext{
		vm:          vm,
		stack:       stack,
		memory:      memory,
		storage:     storage,
		environment: context,
		closureExec: closureExecute,
	}

	is.gasRemaining = evmInt256.FromBigInt(context.Transaction.GasLimit.Int)

	if gasSetting != nil {
		is.gasSetting = gasSetting
	} else {
		is.gasSetting = DefaultGasSetting()
	}

	return is
}
