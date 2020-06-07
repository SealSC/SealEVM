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
	"SealEVM/crypto/hashes"
	"SealEVM/evmErrors"
	"SealEVM/evmInt256"
	"SealEVM/opcodes"
)

func loadMisc() {
	instructionTable[opcodes.SHA3] =  opCodeInstruction{
		action:            sha3Action,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.RETURN] =  opCodeInstruction{
		action:            returnAction,
		requireStackDepth: 2,
		enabled:           true,
		finished:          true,
	}

	instructionTable[opcodes.REVERT] =  opCodeInstruction{
		action:            revertAction,
		requireStackDepth: 2,
		enabled:           true,
		finished:          true,
		returns:           true,
	}

	instructionTable[opcodes.SELFDESTRUCT] =  opCodeInstruction{
		action:            selfDestructAction,
		requireStackDepth: 1,
		enabled:           true,
		finished:          true,
		isWriter:          true,
	}
}

func sha3Action(ctx *instructionsContext) ([]byte, error) {
	mOffset := ctx.stack.Pop()
	mLen := ctx.stack.Pop()

	//gas check
	offset, size, gasLeft, err := ctx.memoryGasCostAndMalloc(mOffset, mLen)
	if err != nil {
		return nil, err
	}

	bytesCost := ctx.gasSetting.DynamicCost.SHA3ByteCost * size
	if gasLeft < bytesCost {
		return nil, evmErrors.OutOfGas
	} else {
		gasLeft -= bytesCost
	}
	ctx.gasRemaining.SetUint64(gasLeft)

	bytes, err := ctx.memory.Copy(offset, size)
	if err != nil {
		return nil, err
	}

	hash := hashes.Keccak256(bytes)

	i := evmInt256.New(0)
	i.SetBytes(hash)
	ctx.stack.Push(i)
	return nil, nil
}

func returnAction(ctx *instructionsContext) ([]byte, error) {
	mOffset := ctx.stack.Pop()
	mLen := ctx.stack.Pop()
	//gas check
	offset, size, _, err := ctx.memoryGasCostAndMalloc(mOffset, mLen)
	if err != nil {
		return nil, err
	}

	return ctx.memory.Copy(offset, size)
}

func revertAction(ctx *instructionsContext) ([]byte, error) {
	mOffset := ctx.stack.Pop()
	mLen := ctx.stack.Pop()

	//gas check
	offset, size, _, err := ctx.memoryGasCostAndMalloc(mOffset, mLen)
	if err != nil {
		return nil, err
	}

	return ctx.memory.Copy(offset, size)
}

func selfDestructAction(ctx *instructionsContext) ([]byte, error) {
	addr := ctx.stack.Pop()
	contractAddr := ctx.environment.Contract.Namespace.Clone()
	balance, _ := ctx.storage.Balance(contractAddr)
	ctx.storage.BalanceModify(addr, balance, false)
	ctx.storage.Destruct(contractAddr)
	return nil, nil
}
