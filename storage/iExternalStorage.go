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
)

type IExternalStorage interface {
	GetBalance(address *evmInt256.Int) (*evmInt256.Int, error)
	GetContract(address *evmInt256.Int) (*Contract, error)
	GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error)

	HashOfCode(code []byte) *evmInt256.Int
	CreateAddress(caller *evmInt256.Int, tx environment.Transaction) *evmInt256.Int
	CreateFixedAddress(caller *evmInt256.Int, salt *evmInt256.Int, code []byte, tx environment.Transaction) *evmInt256.Int

	CanTransfer(from *evmInt256.Int, to *evmInt256.Int, amount *evmInt256.Int) bool
	Load(n *evmInt256.Int, k *evmInt256.Int) (*evmInt256.Int, error)
}
