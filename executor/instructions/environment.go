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
	"SealEVM/common"
	"SealEVM/evmErrors"
	"SealEVM/evmInt256"
	"SealEVM/opcodes"
	"math/big"
)

func loadEnvironment() {
	instructionTable[opcodes.ADDRESS] = opCodeInstruction{
		doAction: addressAction,
		minStackDepth: 0,
		enabled: true,
	}

	instructionTable[opcodes.BALANCE] = opCodeInstruction{
		doAction: balanceAction,
		minStackDepth: 1,
		enabled: true,
	}

	instructionTable[opcodes.ORIGIN] = opCodeInstruction{
		doAction: originAction,
		minStackDepth: 0,
		enabled: true,
	}
	
	instructionTable[opcodes.CALLER] = opCodeInstruction{
		doAction: callerAction,
		minStackDepth: 0,
		enabled: true,
	}

	instructionTable[opcodes.CALLVALUE] = opCodeInstruction{
		doAction: callValueAction,
		minStackDepth: 0,
		enabled: true,
	}

	instructionTable[opcodes.CALLDATALOAD] = opCodeInstruction{
		doAction: callDataLoadAction,
		minStackDepth: 1,
		enabled: true,
	}

	instructionTable[opcodes.CALLDATASIZE] = opCodeInstruction{
		doAction: callDataSizeAction,
		minStackDepth: 0,
		enabled: true,
	}

	instructionTable[opcodes.CALLDATACOPY] = opCodeInstruction{
		doAction: callDataCopyAction,
		minStackDepth: 3,
		enabled: true,
	}

	instructionTable[opcodes.CODESIZE] = opCodeInstruction{
		doAction: codeSizeAction,
		minStackDepth: 0,
		enabled: true,
	}

	instructionTable[opcodes.CODECOPY] = opCodeInstruction{
		doAction: codeCopyAction,
		minStackDepth: 3,
		enabled: true,
	}

	instructionTable[opcodes.GASPRICE] = opCodeInstruction{
		doAction: gasPriceAction,
		minStackDepth: 0,
		enabled: true,
	}

	instructionTable[opcodes.EXTCODESIZE] = opCodeInstruction{
		doAction: extCodeSizeAction,
		minStackDepth: 1,
		enabled: true,
	}

	instructionTable[opcodes.EXTCODECOPY] = opCodeInstruction{
		doAction: extCodeCopyAction,
		minStackDepth: 4,
		enabled: true,
	}

	instructionTable[opcodes.RETURNDATASIZE] = opCodeInstruction{
		doAction: returnDataSizeAction,
		minStackDepth: 0,
		enabled: true,
	}

	instructionTable[opcodes.RETURNDATACOPY] = opCodeInstruction{
		doAction: returnDataCopyAction,
		minStackDepth: 3,
		enabled: true,
	}

	instructionTable[opcodes.EXTCODEHASH] = opCodeInstruction{
		doAction: extCodeHashAction,
		minStackDepth: 1,
		enabled: true,
	}

	instructionTable[opcodes.BLOCKHASH] = opCodeInstruction{
		doAction: blockHashAction,
		minStackDepth: 1,
		enabled: true,
	}

	instructionTable[opcodes.COINBASE] = opCodeInstruction{
		doAction: coinbaseAction,
		minStackDepth: 0,
		enabled: true,
	}

	instructionTable[opcodes.TIMESTAMP] = opCodeInstruction{
		doAction: timestampAction,
		minStackDepth: 0,
		enabled: true,
	}

	instructionTable[opcodes.NUMBER] = opCodeInstruction{
		doAction: numberAction,
		minStackDepth: 0,
		enabled: true,
	}

	instructionTable[opcodes.DIFFICULTY] = opCodeInstruction{
		doAction: difficultyAction,
		minStackDepth: 0,
		enabled: true,
	}

	instructionTable[opcodes.GASLIMIT] = opCodeInstruction{
		doAction: gasLimitAction,
		minStackDepth: 0,
		enabled: true,
	}

}

func addressAction(ctx *instructionsContext) ([]byte, error) {
	err := ctx.stack.Push(ctx.environment.Contract.Address)
	return nil, err
}

func balanceAction(ctx *instructionsContext) ([]byte, error) {
	addr := ctx.stack.Peek()
	balance, err := ctx.storage.Balance(addr)
	if err != nil {
		return nil, err
	}

	addr.Set(balance.Int)
	return nil, nil
}

func originAction(ctx *instructionsContext) ([]byte, error) {
	err := ctx.stack.Push(ctx.environment.Contract.Address)
	return nil, err
}

func callerAction(ctx *instructionsContext) ([]byte, error) {
	err := ctx.stack.Push(ctx.environment.Message.Caller)
	return nil, err
}

func callValueAction(ctx *instructionsContext) ([]byte, error) {
	err := ctx.stack.Push(ctx.environment.Message.Value)
	return nil, err
}

func callDataLoadAction(ctx *instructionsContext) ([]byte, error) {
	i := ctx.stack.Peek()
	data := common.GetDataFrom(ctx.environment.Message.Data, i.Uint64(), 32)

	i.SetBytes(data)
	return nil, nil
}

func callDataSizeAction(ctx *instructionsContext) ([]byte, error) {
	i := ctx.stack.Peek()
	s := ctx.environment.Message.DataSize()

	i.SetUint64(s)
	return nil, nil
}

func callDataCopyAction(ctx *instructionsContext) ([]byte, error) {
	mOffset, _ := ctx.stack.Pop()
	dOffset, _ := ctx.stack.Pop()
	size,_ := ctx.stack.Pop()

	data := common.GetDataFrom(ctx.environment.Message.Data, dOffset.Uint64(), size.Uint64())
	err := ctx.memory.Store(mOffset.Uint64(), data)
	return nil, err
}

func codeSizeAction(ctx *instructionsContext) ([]byte, error) {
	s := evmInt256.New(int64(len(ctx.environment.Contract.Code)))
	err := ctx.stack.Push(s)
	return nil, err
}

func codeCopyAction(ctx *instructionsContext) ([]byte, error) {
	mOffset, _ := ctx.stack.Pop()
	dOffset, _ := ctx.stack.Pop()
	size,_ := ctx.stack.Pop()

	data := common.GetDataFrom(ctx.environment.Contract.Code, dOffset.Uint64(), size.Uint64())
	err := ctx.memory.Store(mOffset.Uint64(), data)
	return nil, err
}

func gasPriceAction(ctx *instructionsContext) ([]byte, error) {
	err := ctx.stack.Push(ctx.environment.Transaction.GasPrice)
	return nil, err
}

func extCodeSizeAction(ctx *instructionsContext) ([]byte, error) {
	addr := ctx.stack.Peek()
	s, err := ctx.storage.GetCodeSize(addr)
	if err != nil {
		return nil, err
	}
	err = ctx.stack.Push(s)
	return nil, err
}

func extCodeCopyAction(ctx *instructionsContext) ([]byte, error) {
	addr, _ := ctx.stack.Pop()
	mOffset, _ := ctx.stack.Pop()
	dOffset, _ := ctx.stack.Pop()
	size,_ := ctx.stack.Pop()

	code, err := ctx.storage.GetCode(addr)
	if err != nil {
		return nil, err
	}

	data := common.GetDataFrom(code, dOffset.Uint64(), size.Uint64())
	err = ctx.memory.Store(mOffset.Uint64(), data)
	return nil, err
}

func returnDataSizeAction(ctx *instructionsContext) ([]byte, error) {
	err := ctx.stack.Push(evmInt256.New(int64(len(ctx.lastReturn))))
	return nil, err
}

func returnDataCopyAction(ctx *instructionsContext) ([]byte, error) {
	mOffset, _ := ctx.stack.Pop()
	dOffset, _ := ctx.stack.Pop()
	size,_ := ctx.stack.Pop()

	end := big.NewInt(0).Add(dOffset.Int, size.Int)
	if !end.IsUint64() || end.Uint64() > uint64(len(ctx.lastReturn)) {
		return nil, evmErrors.ReturnDataCopyOutOfBounds
	}

	err := ctx.memory.Store(mOffset.Uint64(), ctx.lastReturn[dOffset.Uint64() : end.Uint64()])
	return nil, err
}

func extCodeHashAction(ctx *instructionsContext) ([]byte, error) {
	addr:= ctx.stack.Peek()
	codeHash, err := ctx.storage.GetCodeHash(addr)
	if err != nil {
		return nil, err
	}

	addr.Set(codeHash.Int)
	return nil, nil
}

func blockHashAction(ctx *instructionsContext) ([]byte, error) {
	addr:= ctx.stack.Peek()
	codeHash, err := ctx.storage.GetBlockHash(addr)
	if err != nil {
		return nil, err
	}

	addr.Set(codeHash.Int)
	return nil, nil
}

func coinbaseAction(ctx *instructionsContext) ([]byte, error) {
	err := ctx.stack.Push(ctx.environment.Block.Coinbase)
	return nil, err
}

func timestampAction(ctx *instructionsContext) ([]byte, error) {
	err := ctx.stack.Push(ctx.environment.Block.Timestamp)
	return nil, err
}

func numberAction(ctx *instructionsContext) ([]byte, error) {
	err := ctx.stack.Push(ctx.environment.Block.Number)
	return nil, err
}

func difficultyAction(ctx *instructionsContext) ([]byte, error) {
	err := ctx.stack.Push(ctx.environment.Block.Difficulty)
	return nil, err
}

func gasLimitAction(ctx *instructionsContext) ([]byte, error) {
	err := ctx.stack.Push(ctx.environment.Block.GasLimit)
	return nil, err
}
