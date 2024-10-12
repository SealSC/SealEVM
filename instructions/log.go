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
	"github.com/SealSC/SealEVM/opcodes"
	"github.com/SealSC/SealEVM/types"
)

func loadLog() {
	for i := opcodes.LOG0; i <= opcodes.LOG4; i++ {
		topicCount := int(i - opcodes.LOG0)
		instructionTable[i] = opCodeInstruction{
			action: func(ctx *instructionsContext) (bytes []byte, err error) {
				mOffset := ctx.stack.Pop()
				lSize := ctx.stack.Pop()

				log := &types.Log{
					Address: ctx.environment.Address(),
					Topics:  []types.Topic{},
					Data:    []byte{},
				}

				for t := 0; t < topicCount; t++ {
					log.Topics = append(log.Topics, types.Int256ToTopic(ctx.stack.Pop()))
				}

				log.Data, err = ctx.memory.Copy(mOffset.Uint64(), lSize.Uint64())
				if err != nil {
					return
				}

				ctx.storage.Log(log)
				return nil, nil
			},

			requireStackDepth: topicCount + 2,
			enabled:           true,
			isWriter:          true,
		}
	}
}
