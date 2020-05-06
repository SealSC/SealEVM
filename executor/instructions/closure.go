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
	"SealEVM/evmInt256"
	"SealEVM/opcodes"
)

type ClosureExecute func(ClosureParam) ([]byte, error)

type ClosureParam struct {
	OpCode          opcodes.OpCode
	GasRemaining    *evmInt256.Int
	Code            []byte
	Hash            *evmInt256.Int
	Data            []byte
	Value           *evmInt256.Int
	Salt            *evmInt256.Int
}

func loadClosure() {
	instructionTable[opcodes.CALL] = opCodeInstruction {
		doAction: callAction,
		minStackDepth: 7,
		enabled: true,
	}

	instructionTable[opcodes.CALLCODE] = opCodeInstruction {
		doAction: callCodeAction,
		minStackDepth: 7,
		enabled: true,
	}

	instructionTable[opcodes.DELEGATECALL] = opCodeInstruction {
		doAction: delegateCallAction,
		minStackDepth: 6,
		enabled: true,
	}

	instructionTable[opcodes.DELEGATECALL] = opCodeInstruction {
		doAction: staticCallAction,
		minStackDepth: 6,
		enabled: true,
	}

	instructionTable[opcodes.CREATE] = opCodeInstruction {
		doAction: createAction,
		minStackDepth: 3,
		enabled: true,
	}

	instructionTable[opcodes.CREATE2] = opCodeInstruction {
		doAction: create2Action,
		minStackDepth: 2,
		enabled: true,
	}
}

func commonCall(ctx *instructionsContext, opCode opcodes.OpCode) ([]byte, error) {
	_, _ = ctx.stack.Pop()
	addr, _ := ctx.stack.Pop()
	var v *evmInt256.Int = nil
	if opCode != opcodes.DELEGATECALL && opCode != opcodes.STATICCALL {
		v, _ = ctx.stack.Pop()
	}
	dOffset, _ := ctx.stack.Pop()
	dLen, _ := ctx.stack.Pop()
	rOffset, _ := ctx.stack.Pop()
	rLen, _ := ctx.stack.Pop()

	data, err := ctx.memory.Copy(dOffset.Uint64(), dLen.Uint64())
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

	cParam := ClosureParam {
		OpCode:       opCode,
		GasRemaining: ctx.gasRemaining,
		Code:         contractCode,
		Hash:         contractCodeHash,
		Data:         data,
		Value:        v,
	}

	callRet, err := ctx.closureExec(cParam)
	if err != nil {
		return nil, err
	}

	err = ctx.memory.StoreNBytes(rOffset.Uint64(), rLen.Uint64(), callRet)

	return callRet, err
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
	v, _ := ctx.stack.Pop()
	mOffset, _ := ctx.stack.Pop()
	mSize, _ := ctx.stack.Pop()
	var salt *evmInt256.Int = nil
	if opCode == opcodes.CREATE2 {
		salt, _ = ctx.stack.Pop()
	}

	code, err := ctx.memory.Copy(mOffset.Uint64(), mSize.Uint64())
	if err != nil {
		return nil, err
	}

	cParam := ClosureParam {
		OpCode:         opCode,
		GasRemaining:   ctx.gasRemaining,
		Code:           code,
		Data:           code,
		Value:          v,
		Salt:           salt,
	}

	ret, err := ctx.closureExec(cParam)
	if err != nil {
		_ = ctx.stack.Push(evmInt256.New(0))
	} else {
		_ = ctx.stack.Push(evmInt256.New(1))
	}
	return ret, err
}

func createAction(ctx *instructionsContext) ([]byte, error) {
	return commonCreate(ctx, opcodes.CREATE)
}

func create2Action(ctx *instructionsContext) ([]byte, error) {
	return commonCreate(ctx, opcodes.CREATE2)
}
