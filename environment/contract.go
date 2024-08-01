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

package environment

import (
	"github.com/SealSC/SealEVM/evmErrors"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/opcodes"
)

type Contract struct {
	Namespace *evmInt256.Int
	Code      []byte
	Hash      *evmInt256.Int

	codeDataFlag map[uint64]bool
}

func (c *Contract) IsValidJump(dest uint64) (bool, error) {
	codeLen := uint64(len(c.Code))

	if dest > codeLen {
		return false, evmErrors.JumpOutOfBounds
	}

	if c.codeDataFlag == nil {
		c.markCodeData()
	}

	if c.Code[dest] != byte(opcodes.JUMPDEST) {
		return false, evmErrors.InvalidJumpDest
	}

	if c.codeDataFlag[dest] {
		return false, evmErrors.JumpToNoneOpCode
	}

	return true, nil
}

func (c *Contract) GetOpCode(pc uint64) byte {
	if pc < uint64(len(c.Code)) {
		return c.Code[pc]
	}

	return byte(opcodes.STOP)
}

func (c *Contract) markCodeData() {
	c.codeDataFlag = map[uint64]bool{}
	codeLen := len(c.Code)
	for i := 0; i < codeLen; i++ {
		code := opcodes.OpCode(c.Code[i])
		if code >= opcodes.PUSH1 && code <= opcodes.PUSH32 {
			bytesCnt := int(code - opcodes.PUSH1 + 1)
			nextPC := i + bytesCnt
			for ; bytesCnt > 0; bytesCnt-- {
				c.codeDataFlag[uint64(i+bytesCnt)] = true
			}

			i = nextPC
		}
	}
}
