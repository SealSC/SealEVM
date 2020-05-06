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
	"SealEVM/evmInt256"
	"SealEVM/opcodes"
)

func loadMisc() {
	instructionTable[opcodes.SHA3] =  opCodeInstruction {
		doAction: sha3Action,
		minStackDepth: 3,
		enabled: true,
	}

	instructionTable[opcodes.RETURN] =  opCodeInstruction {
		doAction: returnAction,
		minStackDepth: 2,
		enabled: true,
	}

	instructionTable[opcodes.REVERT] =  opCodeInstruction {
		doAction: revertAction,
		minStackDepth: 2,
		enabled: true,
	}

	instructionTable[opcodes.SELFDESTRUCT] =  opCodeInstruction {
		doAction: selfDestructAction,
		minStackDepth: 1,
		enabled: true,
	}
}

func sha3Action(ctx *instructionsContext) ([]byte, error) {
	mOffset, _ := ctx.stack.Pop()
	mLen, _ := ctx.stack.Pop()
	bytes, err := ctx.memory.Copy(mOffset.Uint64(), mLen.Uint64())
	if err != nil {
		return nil, err
	}

	hash := hashes.Keccak256(bytes)

	i := evmInt256.New(0)
	i.SetBytes(hash)
	err = ctx.stack.Push(i)
	return nil, err
}

func returnAction(ctx *instructionsContext) ([]byte, error) {
	mOffset, _ := ctx.stack.Pop()
	mLen, _ := ctx.stack.Pop()

	return ctx.memory.Copy(mOffset.Uint64(), mLen.Uint64())
}

func revertAction(ctx *instructionsContext) ([]byte, error) {
	mOffset, _ := ctx.stack.Pop()
	mLen, _ := ctx.stack.Pop()

	return ctx.memory.Copy(mOffset.Uint64(), mLen.Uint64())
}

func selfDestructAction(ctx *instructionsContext) ([]byte, error) {
	addr, _ := ctx.stack.Pop()
	contractAddr := ctx.environment.Contract.Address
	balance, _ := ctx.storage.Balance(contractAddr)
	ctx.storage.BalanceModify(addr, balance, false)
	ctx.storage.Destruct(contractAddr)
	return nil, nil
}
