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
	"github.com/SealSC/SealEVM/crypto/hashes"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/opcodes"
	"github.com/SealSC/SealEVM/types"
)

func loadMisc() {
	instructionTable[opcodes.SHA3] = opCodeInstruction{
		action:            sha3Action,
		requireStackDepth: 2,
		enabled:           true,
	}

	instructionTable[opcodes.RETURN] = opCodeInstruction{
		action:            returnAction,
		requireStackDepth: 2,
		enabled:           true,
		finished:          true,
	}

	instructionTable[opcodes.REVERT] = opCodeInstruction{
		action:            revertAction,
		requireStackDepth: 2,
		enabled:           true,
		finished:          true,
		returns:           true,
	}

	instructionTable[opcodes.SELFDESTRUCT] = opCodeInstruction{
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

	bytes, err := ctx.memory.Copy(mOffset.Uint64(), mLen.Uint64())
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

	return ctx.memory.Copy(mOffset.Uint64(), mLen.Uint64())
}

func revertAction(ctx *instructionsContext) ([]byte, error) {
	mOffset := ctx.stack.Pop()
	mLen := ctx.stack.Pop()

	return ctx.memory.Copy(mOffset.Uint64(), mLen.Uint64())
}

func selfDestructAction(ctx *instructionsContext) ([]byte, error) {
	receiver := types.Int256ToAddress(ctx.stack.Pop())
	acc := ctx.environment.Account()

	balance, _ := ctx.storage.Balance(acc.Address)

	if !balance.IsZero() {
		if acc.Address != receiver {
			err := ctx.storage.Transfer(acc.Address, receiver, balance)
			if err != nil {
				return nil, err
			}
		}
	}

	ctx.storage.Destruct(acc.Address)
	return nil, nil
}
