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
	"bytes"
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmErrors"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/types"
)

type TypeOfStorage int

const (
	SStorage TypeOfStorage = 1
	TStorage TypeOfStorage = 2
)

type Storage struct {
	ResultCache     ResultCache
	readOnlyCache   readOnlyCache
	externalStorage IExternalStorage
}

func New(extStorage IExternalStorage) *Storage {
	s := &Storage{
		ResultCache:     NewResultCache(),
		externalStorage: extStorage,
		readOnlyCache: readOnlyCache{
			Contracts: ContractCache{},
			BlockHash: Cache{},
		},
	}

	return s
}

func (s *Storage) Clone() *Storage {
	replica := &Storage{
		ResultCache:     s.ResultCache.Clone(),
		readOnlyCache:   s.readOnlyCache,
		externalStorage: s.externalStorage,
	}

	return replica
}

func (s *Storage) XLoad(address types.Address, slot types.Slot, t TypeOfStorage) (*evmInt256.Int, error) {
	if s.ResultCache.OriginalData == nil || s.ResultCache.CachedData == nil || s.externalStorage == nil {
		return nil, evmErrors.StorageNotInitialized
	}

	if t != SStorage && t != TStorage {
		return nil, evmErrors.InvalidTypeOfStorage()
	}

	var err error = nil
	i := s.ResultCache.XCachedLoad(address, slot, t)
	if i == nil {
		if t == SStorage {
			i, err = s.externalStorage.Load(address, slot)
		} else {
			i = evmInt256.New(0)
		}

		if err != nil {
			return nil, evmErrors.NoSuchDataInTheStorage(err)
		}

		s.ResultCache.XCachedStore(address, slot, i, t)
		s.ResultCache.XOriginalStore(address, slot, i, t)
	}

	return i, nil
}

func (s *Storage) XStore(address types.Address, slot types.Slot, val *evmInt256.Int, t TypeOfStorage) {
	s.ResultCache.XCachedStore(address, slot, val, t)
}

func (s *Storage) CanTransfer(from types.Address, to types.Address, amount *evmInt256.Int) bool {
	balance, err := s.Balance(from)
	if err != nil {
		return false
	}

	return balance.Cmp(amount.Int) >= 0
}

func (s *Storage) BalanceModify(address types.Address, value *evmInt256.Int, neg bool) {
	s.Balance(address)

	b, exists := s.ResultCache.Balance[address]
	if !exists {
		b = &balance{
			Address: address,
			Balance: evmInt256.New(0),
		}
		s.ResultCache.Balance[address] = b
	}

	if neg {
		b.Balance.Int.Sub(b.Balance.Int, value.Int)
	} else {
		b.Balance.Int.Add(b.Balance.Int, value.Int)
	}
}

func (s *Storage) Log(log *types.Log) {
	*s.ResultCache.Logs = append(*s.ResultCache.Logs, log)

	return
}

func (s *Storage) Destruct(address types.Address) {
	s.ResultCache.Destructs[address] = address
}

type commonGetterFunc func(types.Slot) (*evmInt256.Int, error)

func (s *Storage) commonGetter(slot types.Slot, cache Cache, getterFunc commonGetterFunc) (*evmInt256.Int, error) {
	if b, exists := cache[slot]; exists {
		return evmInt256.FromBigInt(b.Int), nil
	}

	b, err := getterFunc(slot)
	if err == nil {
		cache[slot] = b
	}

	return b, err
}

func (s *Storage) Balance(address types.Address) (*evmInt256.Int, error) {
	b, exist := s.ResultCache.Balance[address]
	if exist {
		return b.Balance.Clone(), nil
	}

	ba, err := s.externalStorage.GetBalance(address)
	if err != nil {
		b = &balance{
			Address: address,
			Balance: evmInt256.New(0),
		}
	} else {
		b = &balance{
			Address: address,
			Balance: ba,
		}
	}

	s.ResultCache.Balance[address] = b

	return b.Balance.Clone(), nil
}

func (s *Storage) getContract(address types.Address) (*environment.Contract, error) {
	contract := s.ResultCache.NewContracts.Get(address)
	if contract != nil {
		return contract, nil
	}

	contract = s.readOnlyCache.Contracts.Get(address)
	if contract != nil {
		return contract, nil
	}

	contract, err := s.externalStorage.GetContract(address)
	if err != nil {
		return nil, err
	}

	if contract == nil {
		return nil, evmErrors.InvalidExternalStorageResult
	}

	s.readOnlyCache.Contracts.Set(contract)

	return contract, nil
}

func (s *Storage) GetCode(address types.Address) ([]byte, error) {
	contract, err := s.getContract(address)
	if err != nil {
		return nil, err
	}

	return contract.Code, err
}

func (s *Storage) GetCodeSize(address types.Address) (*evmInt256.Int, error) {
	contract, err := s.getContract(address)
	if err != nil {
		return nil, err
	}

	return evmInt256.New(int64(contract.CodeSize)), err
}

func (s *Storage) HashOfCode(code []byte) types.Hash {
	return s.externalStorage.HashOfCode(code)
}

func (s *Storage) GetCodeHash(address types.Address) (*types.Hash, error) {
	contract, err := s.getContract(address)
	if err != nil {
		return nil, err
	}

	return &contract.CodeHash, err
}

func (s *Storage) GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error) {
	var slot types.Slot
	slot.SetBytes(block.Bytes())
	if hash, exists := s.readOnlyCache.BlockHash[slot]; exists {
		return hash, nil
	}

	hash, err := s.externalStorage.GetBlockHash(block)
	if err == nil {
		s.readOnlyCache.BlockHash[slot] = hash
	}

	return hash, err
}

func (s *Storage) NewContract(address types.Address, code []byte) {
	s.ResultCache.NewContracts.Set(&environment.Contract{
		Address:  address,
		Code:     bytes.Clone(code),
		CodeHash: s.externalStorage.HashOfCode(code),
		CodeSize: uint64(len(code)),
	})
}

func (s *Storage) CreateAddress(caller types.Address, tx environment.Transaction) types.Address {
	return s.externalStorage.CreateAddress(caller, tx)
}

func (s *Storage) CreateFixedAddress(caller types.Address, salt types.Hash, code []byte, tx environment.Transaction) types.Address {
	return s.externalStorage.CreateFixedAddress(caller, salt, code, tx)
}

func (s *Storage) GetExternalStorage() IExternalStorage {
	return s.externalStorage
}

func (s *Storage) ClearCache() {
	s.ResultCache = NewResultCache()
}
