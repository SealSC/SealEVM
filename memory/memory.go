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

package memory

import "SealEVM/evmErrors"

type memory struct {
	cell []byte
}

func New() *memory {
	return &memory{}
}

func (m *memory)Malloc(offset uint64, size uint64) []byte {
	mLen := uint64(len(m.cell))
	bound := offset + size
	if mLen < bound {
		newMem := make([]byte, bound - mLen)
		m.cell = append(m.cell, newMem...)
	}

	return m.cell[offset : bound]
}

func (m *memory) Load(offset uint64, length uint64) ([]byte, error) {
	if offset + length > uint64(len(m.cell)) {
		return nil, evmErrors.OutOfMemory
	}

	return m.cell[offset : offset + length], nil
}

func (m *memory) Store(offset uint64, data []byte) error {
	dLen := uint64(len(data))
	if dLen + offset > uint64(len(m.cell)) {
		return evmErrors.OutOfMemory
	}

	copy(m.cell[offset : offset + dLen], data)
	return nil
}

func (m *memory) Copy(offset uint64, length uint64) ([]byte, error) {
	if offset + length > uint64(len(m.cell)) {
		return nil, evmErrors.OutOfMemory
	}

	var ret []byte

	if length == 0 {
		return ret, nil
	}

	copy(ret, m.cell[offset : offset + length])
	return ret, nil
}

func (m *memory) All() []byte {
	return m.cell
}

