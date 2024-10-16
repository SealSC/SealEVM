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
	"github.com/SealSC/SealEVM/opcodes"
	"github.com/SealSC/SealEVM/types"
)

type ClosureExecute func(ClosureParam) ([]byte, error)

type ClosureParam struct {
	VM       interface{}
	OpCode   opcodes.OpCode
	GasLimit *evmInt256.Int

	Called  types.Address
	Message *environment.Message
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
	_ = ctx.stack.Pop() //gas was calculated before execute.
	addr := types.Int256ToAddress(ctx.stack.Pop())
	v := evmInt256.New(0)

	caller := ctx.environment.Address()
	if opCode == opcodes.CALL || opCode == opcodes.CALLCODE {
		v = ctx.stack.Pop()
	} else if opCode == opcodes.DELEGATECALL {
		v = ctx.environment.Message.Value
		caller = ctx.environment.Message.Caller
	}

	dOffset := ctx.stack.Pop()
	dLen := ctx.stack.Pop()
	rOffset := ctx.stack.Pop()
	rLen := ctx.stack.Pop()

	data, err := ctx.memory.Copy(dOffset.Uint64(), dLen.Uint64())
	if err != nil {
		return nil, err
	}

	cParam := ClosureParam{
		VM:       ctx.vm,
		OpCode:   opCode,
		GasLimit: evmInt256.New(ctx.callGasLimit),
		Called:   addr,
		Message: &environment.Message{
			Caller: caller,
			Value:  v,
			Data:   data,
		},
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

	err = ctx.memory.StoreNBytes(rOffset.Uint64(), rLen.Uint64(), callRet)
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
	salt := evmInt256.New(0)
	if opCode == opcodes.CREATE2 {
		salt = ctx.stack.Pop()
	}

	code, err := ctx.memory.Copy(mOffset.Uint64(), mSize.Uint64())
	if err != nil {
		return nil, err
	}

	var addr types.Address
	var caller = ctx.environment.Address()

	if opcodes.CREATE == opCode {
		addr = ctx.storage.CreateAddress(caller, ctx.environment.Transaction)
	} else {
		addr = ctx.storage.CreateFixedAddress(caller, types.Int256ToHash(salt), code, ctx.environment.Transaction)
	}

	var ret []byte
	newCA, err := ctx.storage.AccountWithoutCache(addr)

	if err == nil {
		newCA.Contract = &environment.Contract{
			Code:     code,
			CodeHash: types.Hash{},
			CodeSize: uint64(len(code)),
		}

		ctx.storage.CacheAccount(newCA, true)

		cParam := ClosureParam{
			VM:       ctx.vm,
			OpCode:   opCode,
			GasLimit: ctx.gasRemaining.Clone(),
			Called:   addr,
			Message: &environment.Message{
				Caller: caller,
				Value:  v,
				Data:   nil,
			},
		}

		ret, err = ctx.closureExec(cParam)
	}

	if err != nil {
		ctx.stack.Push(evmInt256.New(0))
		if err != evmErrors.RevertErr {
			ret = nil
		}

		ctx.storage.RemoveCachedAccount(addr)
	} else {
		ctx.stack.Push(addr.Int256())
		ctx.storage.UpdateAccountContract(addr, ret)
	}

	return ret, nil
}

func createAction(ctx *instructionsContext) ([]byte, error) {
	return commonCreate(ctx, opcodes.CREATE)
}

func create2Action(ctx *instructionsContext) ([]byte, error) {
	return commonCreate(ctx, opcodes.CREATE2)
}
