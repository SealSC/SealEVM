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
	"github.com/SealSC/SealEVM/gasSetting"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/opcodes"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
)

type instructionsContext struct {
	stack       *stack.Stack
	memory      *memory.Memory
	storage     *storage.Storage
	environment *environment.Context

	vm interface{}

	pc           uint64
	readOnly     bool
	gasSetting   *gasSetting.Setting
	lastReturn   []byte
	gasRemaining *evmInt256.Int
	callGasLimit uint64
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
	RefundGasFormCall(uint64)
	GetGasLeft() uint64
	GetGasSetting() *gasSetting.Setting
	SetReadOnly()
	IsReadOnly() bool
	ExitOpCode() opcodes.OpCode
}

var instructionTable [opcodes.MaxOpCodesCount]opCodeInstruction

func (i *instructionsContext) SetGasLimit(gasLimit uint64) {
	i.gasRemaining.SetUint64(gasLimit)
}

func (i *instructionsContext) RefundGasFormCall(gasLeft uint64) {
	i.gasRemaining.Add(evmInt256.New(gasLeft))
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

func (i *instructionsContext) GetGasSetting() *gasSetting.Setting {
	return i.gasSetting
}

func (i *instructionsContext) ExitOpCode() opcodes.OpCode {
	return i.exitOpCode
}

func (i *instructionsContext) calcGas(code opcodes.OpCode, gasRemaining uint64) (uint64, error) {
	if code == opcodes.CALL || code == opcodes.CALLCODE || code == opcodes.STATICCALL || code == opcodes.DELEGATECALL {
		if callCost := i.gasSetting.CallCost[code]; callCost != nil {
			if gasRemaining < 100 {
				return 0, evmErrors.OutOfGas
			}

			memExp, gasCost, sendGas, err := callCost(code, gasRemaining, i.stack, i.memory, i.storage)
			if err != nil {
				return gasRemaining, err
			}

			if gasRemaining < gasCost {
				return 0, evmErrors.OutOfGas
			}

			i.callGasLimit = sendGas
			gasRemaining -= gasCost
			i.memory.Malloc(memExp)
		}

		return gasRemaining, nil
	}

	if dynamicCost := i.gasSetting.CommonDynamicCost[code]; dynamicCost != nil {
		memExp, gasCost, err := dynamicCost(i.environment.Contract, i.stack, i.memory, i.storage)
		if err != nil {
			return gasRemaining, err
		}

		i.memory.Malloc(memExp)
		if gasRemaining < gasCost {
			return gasRemaining, evmErrors.OutOfGas
		}

		gasRemaining -= gasCost

		return gasRemaining, nil
	}

	constCost := i.gasSetting.ConstCost[code]
	if gasRemaining < constCost {
		return 0, evmErrors.OutOfGas
	}

	gasRemaining -= constCost

	return gasRemaining, nil
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

		gasLeft, gasErr := i.calcGas(opcodes.OpCode(opCode), i.gasRemaining.Uint64())
		if gasErr != nil {
			err = gasErr
			break
		}

		i.gasRemaining.SetUint64(gasLeft)

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

	if i.exitOpCode == opcodes.REVERT {
		err = evmErrors.RevertErr
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
	gasCfg *gasSetting.Setting,
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

	if gasCfg != nil {
		is.gasSetting = gasCfg
	} else {
		is.gasSetting = gasSetting.Get()
	}

	return is
}
