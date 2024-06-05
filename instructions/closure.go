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
	"github.com/SealSC/SealEVM/evmErrors"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/opcodes"
)

type ClosureExecute func(ClosureParam) ([]byte, error)

type ClosureParam struct {
	VM           interface{}
	OpCode       opcodes.OpCode
	GasRemaining *evmInt256.Int

	ContractAddress *evmInt256.Int
	ContractHash    *evmInt256.Int
	ContractCode    []byte

	CallData   []byte
	CallValue  *evmInt256.Int
	CreateSalt *evmInt256.Int
}

func loadClosure() {
	instructionTable[opcodes.CALL] = opCodeInstruction{
		action:            callAction,
		requireStackDepth: 7,
		enabled:           true,
		returns:           true,
	}

	instructionTable[opcodes.CALLCODE] = opCodeInstruction{
		action:            callCodeAction,
		requireStackDepth: 7,
		enabled:           true,
		returns:           true,
	}

	instructionTable[opcodes.DELEGATECALL] = opCodeInstruction{
		action:            delegateCallAction,
		requireStackDepth: 6,
		enabled:           true,
		returns:           true,
	}

	instructionTable[opcodes.STATICCALL] = opCodeInstruction{
		action:            staticCallAction,
		requireStackDepth: 6,
		enabled:           true,
		returns:           true,
	}

	instructionTable[opcodes.CREATE] = opCodeInstruction{
		action:            createAction,
		requireStackDepth: 3,
		enabled:           true,
		returns:           true,
		isWriter:          true,
	}

	instructionTable[opcodes.CREATE2] = opCodeInstruction{
		action:            create2Action,
		requireStackDepth: 4,
		enabled:           true,
		returns:           true,
		isWriter:          true,
	}
}

func commonCall(ctx *instructionsContext, opCode opcodes.OpCode) ([]byte, error) {
	_ = ctx.stack.Pop()
	addr := ctx.stack.Pop()
	var v *evmInt256.Int = nil
	if opCode != opcodes.DELEGATECALL && opCode != opcodes.STATICCALL {
		v = ctx.stack.Pop()
	}
	dOffset := ctx.stack.Pop()
	dLen := ctx.stack.Pop()
	rOffset := ctx.stack.Pop()
	rLen := ctx.stack.Pop()

	//gas check
	offset, size, _, err := ctx.memoryGasCostAndMalloc(dOffset, dLen)
	if err != nil {
		return nil, err
	}

	data, err := ctx.memory.Copy(offset, size)
	if err != nil {
		return nil, err
	}

	contractCode, err := ctx.storage.GetCode(addr)
	if err != nil {
		return nil, err
	}

	contractCodeHash, err := ctx.storage.GetCodeHash(addr)
	if err != nil {
		return nil, err
	}

	cParam := ClosureParam{
		VM:              ctx.vm,
		OpCode:          opCode,
		GasRemaining:    ctx.gasRemaining,
		ContractAddress: addr,
		ContractCode:    contractCode,
		ContractHash:    contractCodeHash,
		CallData:        data,
		CallValue:       v,
	}

	callRet, err := ctx.closureExec(cParam)
	if err != nil && err != evmErrors.RevertErr {
		ctx.stack.Push(evmInt256.New(0))
		return callRet, nil
	} else if err == evmErrors.RevertErr {
		ctx.stack.Push(evmInt256.New(0))
	} else {
		ctx.stack.Push(evmInt256.New(1))
	}

	//gas check
	offset, size, _, err = ctx.memoryGasCostAndMalloc(rOffset, rLen)
	if err != nil {
		return nil, err
	}

	err = ctx.memory.StoreNBytes(offset, size, callRet)
	if err != nil {
		return nil, err
	}

	return callRet, nil
}

func callAction(ctx *instructionsContext) ([]byte, error) {
	return commonCall(ctx, opcodes.CALL)
}

func callCodeAction(ctx *instructionsContext) ([]byte, error) {
	return commonCall(ctx, opcodes.CALLCODE)
}

func delegateCallAction(ctx *instructionsContext) ([]byte, error) {
	return commonCall(ctx, opcodes.DELEGATECALL)
}

func staticCallAction(ctx *instructionsContext) ([]byte, error) {
	return commonCall(ctx, opcodes.STATICCALL)
}

func commonCreate(ctx *instructionsContext, opCode opcodes.OpCode) ([]byte, error) {
	v := ctx.stack.Pop()
	mOffset := ctx.stack.Pop()
	mSize := ctx.stack.Pop()
	var salt *evmInt256.Int = nil
	if opCode == opcodes.CREATE2 {
		salt = ctx.stack.Pop()
	}

	//gas check
	offset, size, _, err := ctx.memoryGasCostAndMalloc(mOffset, mSize)
	if err != nil {
		return nil, err
	}

	code, err := ctx.memory.Copy(offset, size)
	if err != nil {
		return nil, err
	}

	cParam := ClosureParam{
		VM:           ctx.vm,
		OpCode:       opCode,
		GasRemaining: ctx.gasRemaining,
		ContractCode: code,
		CallData:     []byte{},
		CallValue:    v,
		CreateSalt:   salt,
	}

	var addr *evmInt256.Int
	if opcodes.CREATE == opCode {
		addr = ctx.storage.ExternalStorage.CreateAddress(ctx.environment.Message.Caller, ctx.environment.Transaction)
	} else {
		addr = ctx.storage.ExternalStorage.CreateFixedAddress(ctx.environment.Message.Caller, salt, ctx.environment.Transaction)
	}

	ret, err := ctx.closureExec(cParam)
	if err != nil {
		ctx.stack.Push(evmInt256.New(0))
		if err != evmErrors.RevertErr {
			ret = nil
		}
	} else {
		ctx.stack.Push(addr)
		_ = ctx.storage.NewContract(addr, ret)
	}

	return ret, nil
}

func createAction(ctx *instructionsContext) ([]byte, error) {
	return commonCreate(ctx, opcodes.CREATE)
}

func create2Action(ctx *instructionsContext) ([]byte, error) {
	return commonCreate(ctx, opcodes.CREATE2)
}
