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

type Message struct {
	Caller  []byte
	Value   []byte
	Data    []byte
}

func (m Message) GetData(offset uint64, size uint64) []byte {
	ret := make([]byte, size, size)
	dLen := uint64(len(m.Data))
	if dLen < offset {
		return ret
	}

	end := offset + size
	if dLen < end {
		end = dLen
	}

	copy(ret, m.Data[offset:end])
	return ret
}

func (m Message) DataSize() uint64 {
	return uint64(len(m.Data))
}
