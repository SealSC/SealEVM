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
	"github.com/SealSC/SealEVM/utils"
)

const (
	expandUnit = 32
)

type Memory struct {
	cell        []byte
	lastGasCost uint64
}

func New() *Memory {
	return &Memory{}
}

func (m *Memory) gasCost(size uint64) uint64 {
	wordSize := utils.ToWordSize(size)
	size = wordSize * 32

	square := wordSize * wordSize
	linCoef := wordSize * 3
	quadCoef := square / 512

	return linCoef + quadCoef
}

func (m *Memory) CalculateMallocSizeAndGas(offset *evmInt256.Int, size *evmInt256.Int) (uint64, uint64, error) {
	mLen := uint64(len(m.cell))
	bound := offset.Clone()
	bound.Add(size).ExtendedAlign(expandUnit)

	if !bound.IsUint64() {
		return 0, 0, evmErrors.OutOfMemory
	}

	boundUint := bound.Uint64()

	var expandSize uint64 = 0
	var gasCost uint64 = 0

	if mLen < boundUint {
		newCost := m.gasCost(boundUint)

		gasCost = newCost - m.lastGasCost
		expandSize = boundUint - mLen
	}

	return expandSize, gasCost, nil
}

func (m *Memory) WillIncrease(offset evmInt256.Int, size evmInt256.Int) (o uint64, s uint64, i uint64, err error) {
	mLen := uint64(len(m.cell))
	bound := evmInt256.FromBigInt(offset.Int)
	bound.Add(&size).ExtendedAlign(expandUnit)

	if !bound.IsUint64() {
		return 0, 0, 0, evmErrors.OutOfMemory
	}

	boundUint := bound.Uint64()

	if mLen < boundUint {
		i = boundUint - mLen
	}

	return offset.Uint64(), size.Uint64(), i, nil
}

func (m *Memory) Malloc(length uint64) {
	if length == 0 {
		return
	}
	newMem := make([]byte, length)
	m.cell = append(m.cell, newMem...)
}

func (m *Memory) Map(offset uint64, length uint64) ([]byte, error) {
	if offset+length > uint64(len(m.cell)) {
		return nil, evmErrors.OutOfMemory
	}

	return m.cell[offset : offset+length], nil
}

func (m *Memory) Store(offset uint64, data []byte) error {
	dLen := uint64(len(data))
	if dLen+offset > uint64(len(m.cell)) {
		return evmErrors.OutOfMemory
	}

	copy(m.cell[offset:offset+dLen], data)
	return nil
}

func (m *Memory) StoreNBytes(offset uint64, n uint64, data []byte) error {
	if offset+n > uint64(len(m.cell)) {
		return evmErrors.OutOfMemory
	}

	copy(m.cell[offset:offset+n], data)
	return nil
}

func (m *Memory) Set(idx uint64, data byte) error {
	if idx > uint64(len(m.cell))-1 {
		return evmErrors.OutOfMemory
	}

	m.cell[idx] = data
	return nil
}

func (m *Memory) Copy(offset uint64, length uint64) ([]byte, error) {
	if offset+length > uint64(len(m.cell)) {
		return nil, evmErrors.OutOfMemory
	}

	ret := make([]byte, length, length)
	if length == 0 {
		return ret, nil
	}

	copy(ret, m.cell[offset:offset+length])
	return ret, nil
}

func (m *Memory) MCopy(dst uint64, src uint64, length uint64) error {
	if length == 0 {
		return nil
	}

	mSize := uint64(len(m.cell))
	if src+length > mSize {
		return evmErrors.OutOfMemory
	}

	if dst+length > mSize {
		return evmErrors.OutOfMemory
	}

	copy(m.cell[dst:], m.cell[src:src+length])
	return nil
}

func (m *Memory) Size() int64 {
	return int64(len(m.cell))
}

func (m *Memory) All() []byte {
	return m.cell
}
