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
	"SealEVM/opcodes"
)

func loadLog() {
	for i := opcodes.LOG0; i <= opcodes.LOG4; i++ {
		topicCount := int(i - opcodes.LOG0)
		instructionTable[i] = opCodeInstruction {
			doAction: func(ctx *instructionsContext) (bytes []byte, err error) {
				mOffset, _ := ctx.stack.Pop()
				lSize, _ := ctx.stack.Pop()
				var topics [][]byte

				for t := 0; t < topicCount; t++ {
					topic, _ := ctx.stack.Pop()
					topicBytes := common.EVMIntToHashBytes(topic)
					topics = append(topics, topicBytes[:])
				}

				log, err := ctx.memory.Copy(mOffset.Uint64(), lSize.Uint64())
				if err != nil {
					return
				}

				ctx.storage.Log(ctx.environment.Contract.Address, topics, log, ctx.environment)
				return nil, nil
			},
		}
	}
}
