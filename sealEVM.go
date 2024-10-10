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
	"github.com/SealSC/SealEVM/gasSetting"
	"github.com/SealSC/SealEVM/instructions"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/opcodes"
	"github.com/SealSC/SealEVM/precompiledContracts"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
	"github.com/SealSC/SealEVM/types"
	"github.com/SealSC/SealEVM/utils"
)

type EVMResultCallback func(result ExecuteResult, err error)
type EVMParam struct {
	MaxStackDepth  int
	ExternalStore  storage.IExternalStorage
	ResultCallback EVMResultCallback
	Context        *environment.Context
	GasSetting     *gasSetting.Setting
}

type EVM struct {
	depth        uint64
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

	evm.instructions = instructions.New(evm, evm.stack, evm.memory, evm.storage, evm.context, param.GasSetting, closure)

	return evm
}

func NewWithCache(param EVMParam, s *storage.Storage) *EVM {
	if param.Context.Block.GasLimit.Cmp(param.Context.Transaction.GasLimit.Int) < 0 {
		param.Context.Transaction.GasLimit = evmInt256.FromBigInt(param.Context.Block.GasLimit.Int)
	}

	evm := &EVM{
		stack:        stack.New(param.MaxStackDepth),
		memory:       memory.New(),
		storage:      s.Clone(),
		context:      param.Context,
		instructions: nil,
		resultNotify: param.ResultCallback,
	}

	evm.instructions = instructions.New(evm, evm.stack, evm.memory, evm.storage, evm.context, param.GasSetting, closure)

	return evm
}

func (e *EVM) subResult(result ExecuteResult, err error) {
	if err == nil && result.ExitOpCode != opcodes.REVERT {
		storage.MergeResultCache(&result.StorageCache, &e.storage.ResultCache)
	}
}

func (e *EVM) executePreCompiled(address types.Address, input []byte) (ExecuteResult, error) {
	addrInt := address.Int256().Uint64()
	contract := precompiledContracts.GetContract(addrInt)
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

func (e *EVM) createContract(env *environment.Context) *environment.Contract {
	return &environment.Contract{
		Address:  e.storage.CreateAddress(env.Message.Caller, env.Transaction),
		Code:     env.Message.Data,
		CodeHash: e.storage.HashOfCode(env.Message.Data),
		CodeSize: uint64(len(env.Message.Data)),
	}
}

func (e *EVM) useIntrinsicGas() (uint64, error) {
	gasLeft := e.instructions.GetGasLeft()
	gasCost := e.instructions.GetGasSetting().IntrinsicCost(e.context.Message.Data, e.context.Contract)
	if gasLeft < gasCost {
		e.instructions.SetGasLimit(0)
		return 0, evmErrors.OutOfGas
	}

	gasLeft -= gasCost
	e.instructions.SetGasLimit(gasLeft)

	return gasLeft, nil
}

func (e *EVM) ExecuteContract() (ExecuteResult, error) {
	var isCreation = false
	var err error
	var gasLeft uint64

	if e.depth == 0 {
		gasLeft, err = e.useIntrinsicGas()
	} else {
		gasLeft = e.instructions.GetGasLeft()
	}

	result := ExecuteResult{
		ResultData:   nil,
		GasLeft:      gasLeft,
		StorageCache: e.storage.ResultCache,
	}

	if err != nil {
		return result, err
	}

	if e.context.Contract == nil {
		isCreation = true
		e.context.Contract = e.createContract(e.context)
	}

	if e.context.Message.Value == nil {
		e.context.Message.Value = evmInt256.New(0)
	}

	contractAddr := e.context.Contract.Address

	//doing transfer when value > 0
	if e.context.Message.Value.Sign() > 0 {
		msg := e.context.Message
		if e.instructions.IsReadOnly() {
			return result, evmErrors.WriteProtection
		}

		if e.storage.CanTransfer(msg.Caller, contractAddr, msg.Value) {
			transErr := e.storage.Transfer(msg.Caller, contractAddr, msg.Value)
			if transErr != nil {
				return result, err
			}
		} else {
			return result, evmErrors.InsufficientBalance
		}
	}

	if precompiledContracts.IsPrecompiledContract(contractAddr) {
		return e.executePreCompiled(contractAddr, e.context.Message.Data)
	}

	execRet, gasLeft, err := e.instructions.ExecuteContract()

	if err == nil {
		if isCreation {
			storeCost, storeErr := e.instructions.GetGasSetting().ContractStoreCost(execRet, gasLeft)
			if storeErr != nil {
				err = storeErr
				gasLeft = 0
			} else {
				gasLeft -= storeCost
			}
		}

	}

	result.GasLeft = gasLeft
	result.ResultData = execRet
	result.ExitOpCode = e.instructions.ExitOpCode()

	if err != nil {
		result.StorageCache = storage.NewResultCache()
		e.storage.ClearCache()
	}

	if e.resultNotify != nil {
		e.resultNotify(result, err)
	}

	return result, err
}

func (e *EVM) getClosureDefaultEVM(param instructions.ClosureParam) *EVM {
	newEVM := NewWithCache(EVMParam{
		MaxStackDepth:  1024,
		ExternalStore:  e.storage.GetExternalStorage(),
		ResultCallback: e.subResult,
		Context: &environment.Context{
			Block:       e.context.Block,
			Transaction: e.context.Transaction,
			Message: environment.Message{
				Data: param.CallData,
			},
		},
		GasSetting: e.instructions.GetGasSetting(),
	}, e.storage)

	newEVM.context.Contract = &environment.Contract{
		Address:  param.ContractAddress,
		Code:     param.ContractCode,
		CodeHash: param.ContractHash,
	}

	newEVM.instructions.SetGasLimit(param.GasLimit.Uint64())

	return newEVM
}

func (e *EVM) commonCall(param instructions.ClosureParam, depth uint64) ([]byte, error) {
	newEVM := e.getClosureDefaultEVM(param)
	newEVM.depth = depth

	//set storage address and call value
	switch param.OpCode {
	case opcodes.CALL:
		newEVM.context.Contract.Address = param.ContractAddress
		newEVM.context.Message.Value = param.CallValue
		newEVM.context.Message.Caller = e.context.Contract.Address
	case opcodes.STATICCALL:
		newEVM.context.Contract.Address = param.ContractAddress
		newEVM.context.Message.Value = param.CallValue
		newEVM.context.Message.Caller = e.context.Contract.Address
	case opcodes.CALLCODE:
		newEVM.context.Contract.Address = e.context.Contract.Address
		newEVM.context.Message.Value = param.CallValue
		newEVM.context.Message.Caller = e.context.Contract.Address

	case opcodes.DELEGATECALL:
		newEVM.context.Contract.Address = e.context.Contract.Address
		newEVM.context.Message.Value = e.context.Message.Value
		newEVM.context.Message.Caller = e.context.Message.Caller
	}

	if param.OpCode == opcodes.STATICCALL || e.instructions.IsReadOnly() {
		newEVM.instructions.SetReadOnly()
	}

	ret, err := newEVM.ExecuteContract()
	if ret.ExitOpCode == opcodes.REVERT {
		err = evmErrors.RevertErr
	}

	e.instructions.RefundGasFormCall(ret.GasLeft)
	return ret.ResultData, err
}

func (e *EVM) commonCreate(param instructions.ClosureParam, depth uint64) ([]byte, error) {
	newEVM := e.getClosureDefaultEVM(param)

	newEVM.depth = depth
	newEVM.context.Contract.Address = param.ContractAddress
	newEVM.context.Message.Value = param.CallValue
	newEVM.context.Message.Caller = e.context.Contract.Address

	ret, err := newEVM.ExecuteContract()
	if ret.ExitOpCode == opcodes.REVERT {
		err = evmErrors.RevertErr
	}
	e.instructions.SetGasLimit(ret.GasLeft)
	return ret.ResultData, err
}

func closure(param instructions.ClosureParam) ([]byte, error) {
	evm, ok := param.VM.(*EVM)
	if !ok {
		return nil, evmErrors.InvalidEVMInstance
	}

	evm.depth += 1
	defer func() {
		evm.depth -= 1
	}()

	if evm.depth > utils.MaxClosureDepth {
		return nil, evmErrors.ClosureDepthOverflow
	}

	switch param.OpCode {
	case opcodes.CALL, opcodes.CALLCODE, opcodes.DELEGATECALL, opcodes.STATICCALL:
		return evm.commonCall(param, evm.depth)
	case opcodes.CREATE, opcodes.CREATE2:
		return evm.commonCreate(param, evm.depth)
	}

	return nil, nil
}
