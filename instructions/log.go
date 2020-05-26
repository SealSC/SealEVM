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
	"SealEVM/opcodes"
)

func loadLog() {
	for i := opcodes.LOG0; i <= opcodes.LOG4; i++ {
		topicCount := int(i - opcodes.LOG0)
		instructionTable[i] = opCodeInstruction{
			action: func(ctx *instructionsContext) (bytes []byte, err error) {
				mOffset := ctx.stack.Pop()
				lSize := ctx.stack.Pop()
				var topics [][]byte

				//gas check
				offset, size, gasLeft, err := ctx.memoryGasCostAndMalloc(mOffset, lSize)
				if err != nil {
					return nil, err
				}

				logByteCost := ctx.gasSetting.DynamicCost.LogByteCost * lSize.Uint64()
				if gasLeft < logByteCost {
					return nil, evmErrors.OutOfGas
				} else {
					gasLeft -= logByteCost
				}
				ctx.gasRemaining.SetUint64(gasLeft)

				for t := 0; t < topicCount; t++ {
					topic := ctx.stack.Pop()
					topicBytes := common.EVMIntToHashBytes(topic)
					topics = append(topics, topicBytes[:])
				}

				log, err := ctx.memory.Copy(offset, size)
				if err != nil {
					return
				}

				ctx.storage.Log(ctx.environment.Contract.Namespace, topics, log, *ctx.environment)
				return nil, nil
			},

			requireStackDepth: topicCount + 2,
			enabled:           true,
			isWriter:          true,
		}
	}
}
