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

import (
	"github.com/SealSC/SealEVM/evmErrors"
	"github.com/SealSC/SealEVM/evmInt256"
)

type Memory struct {
	cell []byte
}

func New() *Memory {
	return &Memory{}
}

func (m *Memory) WillIncrease(offset *evmInt256.Int, size *evmInt256.Int) (o uint64, s uint64, i uint64, err error) {
	mLen := uint64(len(m.cell))
	bound := evmInt256.FromBigInt(offset.Int)
	bound.Add(size)

	if !bound.IsUint64() {
		return 0, 0, 0, evmErrors.OutOfMemory
	}

	boundUint := bound.Uint64()

	if mLen < boundUint {
		i = boundUint - mLen
	}

	return offset.Uint64(), size.Uint64(), i, nil
}

func (m *Memory)Malloc(offset uint64, size uint64) []byte {
	mLen := uint64(len(m.cell))
	bound := offset + size
	if mLen < bound {
		newMem := make([]byte, bound - mLen)
		m.cell = append(m.cell, newMem...)
	}

	return m.cell[offset : bound]
}

func (m *Memory) Map(offset uint64, length uint64) ([]byte, error) {
	if offset + length > uint64(len(m.cell)) {
		return nil, evmErrors.OutOfMemory
	}

	return m.cell[offset : offset + length], nil
}

func (m *Memory) Store(offset uint64, data []byte) error {
	dLen := uint64(len(data))
	if dLen + offset > uint64(len(m.cell)) {
		return evmErrors.OutOfMemory
	}

	copy(m.cell[offset : offset + dLen], data)
	return nil
}

func (m *Memory) StoreNBytes(offset uint64, n uint64, data []byte) error {
	if offset + n > uint64(len(m.cell)) {
		return evmErrors.OutOfMemory
	}

	copy(m.cell[offset : offset + n], data)
	return nil
}

func (m *Memory) Set(idx uint64, data byte) error {
	if idx > uint64(len(m.cell)) - 1 {
		return evmErrors.OutOfMemory
	}

	m.cell[idx] = data
	return nil
}

func (m *Memory) Copy(offset uint64, length uint64) ([]byte, error) {
	if offset + length > uint64(len(m.cell)) {
		return nil, evmErrors.OutOfMemory
	}

	ret := make([]byte, length, length)
	if length == 0 {
		return ret, nil
	}

	copy(ret, m.cell[offset : offset + length])
	return ret, nil
}

func (m *Memory) Size() int64 {
	return int64(len(m.cell))
}

func (m *Memory) All() []byte {
	return m.cell
}

