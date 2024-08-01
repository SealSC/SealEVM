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

package common

import (
	"github.com/SealSC/SealEVM/crypto/hashes"
	"github.com/SealSC/SealEVM/evmInt256"
)

const (
	HashLength      = 32
	MaxClosureDepth = 1024
)

var (
	BlankHash = make([]byte, HashLength, HashLength)
	ZeroHash  = hashes.Keccak256(nil)
)

func EVMIntToHashBytes(i *evmInt256.Int) [HashLength]byte {
	iBytes := i.Bytes()
	iLen := len(iBytes)

	var hash [HashLength]byte
	if iLen > HashLength {
		copy(hash[:], iBytes[iLen-HashLength:])
	} else {
		copy(hash[HashLength-iLen:], iBytes)
	}

	return hash
}

func HashBytesToEVMInt(hash [HashLength]byte) (*evmInt256.Int, error) {

	i := evmInt256.New(0)
	i.SetBytes(hash[:])

	return i, nil
}

func BytesDataToEVMIntHash(data []byte) *evmInt256.Int {
	var hashBytes []byte
	srcLen := len(data)
	if srcLen < HashLength {
		hashBytes = LeftPaddingSlice(data, HashLength)
	} else {
		hashBytes = data[:HashLength]
	}

	i := evmInt256.New(0)
	i.SetBytes(hashBytes)

	return i
}

func GetDataFrom(src []byte, offset uint64, size uint64) []byte {
	ret := make([]byte, size, size)
	dLen := uint64(len(src))
	if dLen < offset {
		return ret
	}

	end := offset + size
	if dLen < end {
		end = dLen
	}

	copy(ret, src[offset:end])
	return ret
}

func LeftPaddingSlice(src []byte, toSize int) []byte {
	sLen := len(src)
	if toSize <= sLen {
		return src
	}

	ret := make([]byte, toSize, toSize)
	copy(ret[toSize-sLen:], src)

	return ret
}

func RightPaddingSlice(src []byte, toSize int) []byte {
	sLen := len(src)
	if toSize <= sLen {
		return src
	}

	ret := make([]byte, toSize, toSize)
	copy(ret, src)

	return ret
}
