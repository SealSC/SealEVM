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
	"SealEVM/executor/instructions"
	"SealEVM/memory"
	"SealEVM/stack"
	"SealEVM/storageCache"
)

type EVMParam struct {
	MaxStackDepth   int
	ExternalStore   storageCache.IExternalStorage
	ResultCallback  storageCache.EVMResultCallback
	Context         environment.Context
}

type EVM struct {
	stack           *stack.Stack
	memory          *memory.Memory
	iStore          storageCache.ICache
	context         environment.Context
	instructions    instructions.IInstructions
}

func New(param EVMParam) *EVM {
	evm := &EVM{
		stack:        stack.New(param.MaxStackDepth),
		memory:       memory.New(),
		iStore:       storageCache.New(param.ExternalStore, param.ResultCallback),
		context:      param.Context,
		instructions: nil,
	}

	evm.instructions = instructions.New(evm.stack, evm.memory, evm.iStore, evm.context)

	return evm
}
