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

package SealEVM

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmErrors"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/instructions"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/opcodes"
	"github.com/SealSC/SealEVM/precompiledContracts"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
)

type EVMResultCallback func(result ExecuteResult, err error)
type EVMParam struct {
	MaxStackDepth  int
	ExternalStore  storage.IExternalStorage
	ResultCallback EVMResultCallback
	Context        *environment.Context
}

type EVM struct {
	stack        *stack.Stack
	memory       *memory.Memory
	storage      *storage.Storage
	context      *environment.Context
	instructions instructions.IInstructions
	resultNotify EVMResultCallback
}

type ExecuteResult struct {
	ResultData   []byte
	GasLeft      uint64
	StorageCache storage.ResultCache
	ExitOpCode   opcodes.OpCode
}

func Load() {
	instructions.Load()
}

func New(param EVMParam) *EVM {
	if param.Context.Block.GasLimit.Cmp(param.Context.Transaction.GasLimit.Int) < 0 {
		param.Context.Transaction.GasLimit = evmInt256.FromBigInt(param.Context.Block.GasLimit.Int)
	}

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

func (e *EVM) subResult(result ExecuteResult, err error) {
	if err == nil && result.ExitOpCode != opcodes.REVERT {
		storage.MergeResultCache(&result.StorageCache, &e.storage.ResultCache)
	}
}

func (e *EVM) executePreCompiled(addr uint64, input []byte) (ExecuteResult, error) {
	contract := precompiledContracts.GetContract(addr)
	gasCost := contract.GasCost(input)
	gasLeft := e.instructions.GetGasLeft()

	result := ExecuteResult{
		ResultData:   nil,
		GasLeft:      gasLeft,
		StorageCache: e.storage.ResultCache,
	}

	if gasLeft < gasCost {
		return result, evmErrors.OutOfGas
	}

	execRet, err := contract.Execute(input)
	gasLeft -= gasCost
	e.instructions.SetGasLimit(gasLeft)
	result.ResultData = execRet
	return result, err
}

func (e *EVM) ExecuteContract(doTransfer bool) (ExecuteResult, error) {
	contractAddr := e.context.Contract.Namespace
	gasLeft := e.instructions.GetGasLeft()

	result := ExecuteResult{
		ResultData:   nil,
		GasLeft:      gasLeft,
		StorageCache: e.storage.ResultCache,
	}

	if doTransfer {
		msg := e.context.Message

		//doing transfer for non-zero value
		if msg.Value.Sign() != 0 {
			if e.instructions.IsReadOnly() {
				return result, evmErrors.WriteProtection
			}

			if e.storage.ExternalStorage.CanTransfer(msg.Caller, contractAddr, msg.Value) {
				e.storage.BalanceModify(msg.Caller, msg.Value, true)
				e.storage.BalanceModify(contractAddr, msg.Value, false)
			} else {
				return result, evmErrors.InsufficientBalance
			}
		}

	}

	if contractAddr != nil {
		//check if is precompiled
		if contractAddr.IsUint64() {
			addr := contractAddr.Uint64()
			if addr < precompiledContracts.PrecompiledContractCount() {
				return e.executePreCompiled(addr, e.context.Message.Data)
			}
		}
	}

	execRet, gasLeft, err := e.instructions.ExecuteContract()
	result.ResultData = execRet
	result.ExitOpCode = e.instructions.ExitOpCode()

	if e.resultNotify != nil {
		e.resultNotify(result, err)
	}

	return result, err
}

func (e *EVM) getClosureDefaultEVM(param instructions.ClosureParam) *EVM {
	newEVM := New(EVMParam{
		MaxStackDepth:  1024,
		ExternalStore:  e.storage.ExternalStorage,
		ResultCallback: e.subResult,
		Context: &environment.Context{
			Block:       e.context.Block,
			Transaction: e.context.Transaction,
			Message: environment.Message{
				Data: param.CallData,
			},
		},
	})

	newEVM.context.Contract = environment.Contract{
		Namespace: param.ContractAddress,
		Code:      param.ContractCode,
		Hash:      param.ContractHash,
	}

	return newEVM
}

func (e *EVM) commonCall(param instructions.ClosureParam) ([]byte, error) {
	newEVM := e.getClosureDefaultEVM(param)

	//set storage namespace and call value
	switch param.OpCode {
	case opcodes.CALL:
		newEVM.context.Contract.Namespace = param.ContractAddress
		newEVM.context.Message.Value = param.CallValue
		newEVM.context.Message.Caller = e.context.Contract.Namespace
	case opcodes.STATICCALL:
		newEVM.context.Contract.Namespace = param.ContractAddress
		newEVM.context.Message.Value = param.CallValue
		newEVM.context.Message.Caller = e.context.Contract.Namespace
	case opcodes.CALLCODE:
		newEVM.context.Contract.Namespace = e.context.Contract.Namespace
		newEVM.context.Message.Value = param.CallValue
		newEVM.context.Message.Caller = e.context.Contract.Namespace

	case opcodes.DELEGATECALL:
		newEVM.context.Contract.Namespace = e.context.Contract.Namespace
		newEVM.context.Message.Value = e.context.Message.Value
		newEVM.context.Message.Caller = e.context.Message.Caller
	}

	if param.OpCode == opcodes.STATICCALL || e.instructions.IsReadOnly() {
		newEVM.instructions.SetReadOnly()
	}

	ret, err := newEVM.ExecuteContract(opcodes.CALL == param.OpCode)

	e.instructions.SetGasLimit(ret.GasLeft)
	return ret.ResultData, err
}

func (e *EVM) commonCreate(param instructions.ClosureParam) ([]byte, error) {
	var addr *evmInt256.Int
	if opcodes.CREATE == param.OpCode {
		addr = e.storage.ExternalStorage.CreateAddress(e.context.Message.Caller, e.context.Transaction)
	} else {
		addr = e.storage.ExternalStorage.CreateFixedAddress(e.context.Message.Caller, param.CreateSalt, e.context.Transaction)
	}

	newEVM := e.getClosureDefaultEVM(param)

	newEVM.context.Contract.Namespace = addr
	newEVM.context.Message.Value = param.CallValue
	newEVM.context.Message.Caller = e.context.Contract.Namespace

	ret, err := newEVM.ExecuteContract(true)
	e.instructions.SetGasLimit(ret.GasLeft)
	return ret.ResultData, err
}

func closure(param instructions.ClosureParam) ([]byte, error) {
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
