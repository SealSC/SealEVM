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

package storage

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/types"
)

type IExternalStorage interface {
	GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error)

	GetAccount(address types.Address) (*environment.Account, error)
	AccountExist(address types.Address) bool
	AccountEmpty(address types.Address) bool

	HashOfCode(code []byte) types.Hash
	CreateAddress(caller types.Address, tx environment.Transaction) types.Address
	CreateFixedAddress(caller types.Address, salt types.Hash, code []byte, tx environment.Transaction) types.Address

	Load(address types.Address, slot types.Slot) (*evmInt256.Int, error)
}

type IExternalDataBlockStorage interface {
	GetDataBlock(address types.Address, slot types.Slot) (types.Bytes, error)
}
