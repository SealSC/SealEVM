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

package executor

import (
	"SealEVM/environment"
	"SealEVM/evmErrors"
	"SealEVM/executor/instructions"
	"SealEVM/memory"
	"SealEVM/opcodes"
	"SealEVM/precompiledContracts"
	"SealEVM/stack"
	"SealEVM/storage"
)

type EVMResultCallback func(contractRet[]byte, gasLeft uint64, result storage.ResultCache, err error)
type EVMParam struct {
	MaxStackDepth  int
	ExternalStore  storage.IExternalStorage
	ResultCallback EVMResultCallback
	Context        environment.Context
}

type EVM struct {
	stack           *stack.Stack
	memory          *memory.Memory
	storage         *storage.Storage
	context         environment.Context
	instructions    instructions.IInstructions
	resultNotify    EVMResultCallback
}

func Load() {
	instructions.Load()
}

func New(param EVMParam) *EVM {
	evm := &EVM{
		stack:        stack.New(param.MaxStackDepth),
		memory:       memory.New(),
		storage:      storage.New(param.ExternalStore),
		context:      param.Context,
		instructions: nil,
		resultNotify: param.ResultCallback,
	}

	evm.instructions = instructions.New(evm, evm.stack, evm.memory, evm.storage, evm.context, nil, closure)

	return evm
}

func (e *EVM) subResult(contractRet []byte, gasLeft uint64, cache storage.ResultCache, err error) {
	if err == nil {
		storage.MergeResultCache(&cache, &e.storage.ResultCache)
	}
}

func (e *EVM) executePreCompiled(addr uint64, input []byte) ([]byte, uint64, error) {
	contract := precompiledContracts.Contracts[addr]
	gasCost := contract.GasCost(input)
	gasLeft := e.instructions.GetGasLeft()

	if gasLeft < gasCost {
		return nil, gasLeft, evmErrors.OutOfGas
	}

	ret, err := contract.Execute(input)
	gasLeft -= gasCost
	e.instructions.SetGasLimit(gasLeft)
	return ret, gasLeft, err
}

func (e *EVM) ExecuteContract() ([]byte, uint64, error) {
	contractAddr := e.context.Contract.Namespace

	if contractAddr != nil {
		//transfer first
		msg := e.context.Message
		if e.storage.ExternalStorage.CanTransfer(msg.Caller, contractAddr, msg.Value) {
			e.storage.BalanceModify(msg.Caller, msg.Value, true)
			e.storage.BalanceModify(contractAddr, msg.Value, false)
		} else {
			return nil, e.instructions.GetGasLeft(), evmErrors.InsufficientBalance
		}

		//check if is precompiled
		if contractAddr.IsUint64() {
			addr := contractAddr.Uint64()
			if addr < precompiledContracts.ContractsMaxAddress {
				return e.executePreCompiled(addr, e.context.Message.Data)
			}
		}
	}

	ret, gasLeft, err := e.instructions.ExecuteContract()
	e.resultNotify(ret, gasLeft, e.storage.ResultCache, err)
	return ret, gasLeft, err
}

func (e *EVM) commonCall(param instructions.ClosureParam) ([]byte, error) {
	newEVM := New(EVMParam {
		MaxStackDepth:  1024,
		ExternalStore:  e.storage.ExternalStorage,
		ResultCallback: e.subResult,
		Context:        environment.Context {
			Block:          e.context.Block,
			Transaction:    e.context.Transaction,
			Message:        environment.Message {
				Data:   param.CallData,
			},
		},
	})

	newEVM.context.Contract = environment.Contract {
		Namespace:  param.ContractAddress,
		Code:       param.ContractCode,
		Hash:       param.ContractHash,
	}

	//set storage namespace and call value
	switch param.OpCode {
	case opcodes.CALLCODE:
		newEVM.context.Contract.Namespace = e.context.Contract.Namespace
		newEVM.context.Message.Value = param.CallValue
		newEVM.context.Message.Caller = e.context.Contract.Namespace

	case opcodes.DELEGATECALL:
		newEVM.context.Contract.Namespace = e.context.Contract.Namespace
		newEVM.context.Message.Value = e.context.Message.Value
		newEVM.context.Message.Caller = e.context.Message.Caller
	}

	ret, gasLeft, err := newEVM.instructions.ExecuteContract()

	e.instructions.SetGasLimit(gasLeft)
	return ret, err
}

func (e *EVM) commonCreate(param instructions.ClosureParam) ([]byte, error) {
	//todo: implement create operation code
	return nil, nil
}

func closure(param instructions.ClosureParam) ([]byte, error){
	evm, ok := param.VM.(*EVM)
	if !ok {
		return nil, evmErrors.InvalidEVMInstance
	}

	switch param.OpCode {
	case opcodes.CALL, opcodes.CALLCODE, opcodes.DELEGATECALL, opcodes.STATICCALL:
		return evm.commonCall(param)
	case opcodes.CREATE, opcodes.CREATE2:
		return evm.commonCreate(param)
	}

	return nil, nil
}
