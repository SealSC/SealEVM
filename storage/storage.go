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

func (s *Storage) XLoad(n *evmInt256.Int, k *evmInt256.Int, t TypeOfStorage) (*evmInt256.Int, error) {
	if s.ResultCache.OriginalData == nil || s.ResultCache.CachedData == nil || s.externalStorage == nil {
		return nil, evmErrors.StorageNotInitialized
	}

	nsStr := n.AsStringKey()
	keyStr := k.AsStringKey()

	if t != SStorage && t != TStorage {
		return nil, evmErrors.InvalidTypeOfStorage()
	}

	var err error = nil
	i := s.ResultCache.XCachedLoad(nsStr, keyStr, t)
	if i == nil {
		if t == SStorage {
			i, err = s.externalStorage.Load(n, k)
		} else {
			i = evmInt256.New(0)
		}

		if err != nil {
			return nil, evmErrors.NoSuchDataInTheStorage(err)
		}

		s.ResultCache.XCachedStore(nsStr, keyStr, i, t)
		s.ResultCache.XOriginalStore(nsStr, keyStr, i, t)
	}

	return i, nil
}

func (s *Storage) XStore(n *evmInt256.Int, k *evmInt256.Int, v *evmInt256.Int, t TypeOfStorage) {
	s.ResultCache.XCachedStore(n.AsStringKey(), k.AsStringKey(), v, t)
}

func (s *Storage) CanTransfer(from *evmInt256.Int, to *evmInt256.Int, amount *evmInt256.Int) bool {
	balance, err := s.Balance(from)
	if err != nil {
		return false
	}

	return balance.Cmp(amount.Int) >= 0
}

func (s *Storage) BalanceModify(address *evmInt256.Int, value *evmInt256.Int, neg bool) {
	kString := address.AsStringKey()
	s.Balance(address)

	b, exists := s.ResultCache.Balance[kString]
	if !exists {
		b = &balance{
			Address: evmInt256.FromBigInt(address.Int),
			Balance: evmInt256.New(0),
		}
		s.ResultCache.Balance[address.AsStringKey()] = b
	}

	if neg {
		b.Balance.Int.Sub(b.Balance.Int, value.Int)
	} else {
		b.Balance.Int.Add(b.Balance.Int, value.Int)
	}
}

func (s *Storage) Log(address *evmInt256.Int, topics [][]byte, data []byte, context environment.Context) {
	var theLog = Log{
		Address: address,
		Topics:  topics,
		Data:    data,
	}

	*s.ResultCache.Logs = append(*s.ResultCache.Logs, theLog)

	return
}

func (s *Storage) Destruct(address *evmInt256.Int) {
	s.ResultCache.Destructs[address.AsStringKey()] = address
}

type commonGetterFunc func(*evmInt256.Int) (*evmInt256.Int, error)

func (s *Storage) commonGetter(key *evmInt256.Int, cache Cache, getterFunc commonGetterFunc) (*evmInt256.Int, error) {
	keyStr := key.AsStringKey()
	if b, exists := cache[keyStr]; exists {
		return evmInt256.FromBigInt(b.Int), nil
	}

	b, err := getterFunc(key)
	if err == nil {
		cache[keyStr] = b
	}

	return b, err
}

func (s *Storage) Balance(address *evmInt256.Int) (*evmInt256.Int, error) {
	b, exist := s.ResultCache.Balance[address.AsStringKey()]
	if exist {
		return b.Balance.Clone(), nil
	}

	ba, err := s.externalStorage.GetBalance(address)
	if err != nil {
		b = &balance{
			Address: evmInt256.FromBigInt(address.Int),
			Balance: evmInt256.New(0),
		}
	} else {
		b = &balance{
			Address: address,
			Balance: ba,
		}
	}

	s.ResultCache.Balance[address.AsStringKey()] = b

	return b.Balance.Clone(), nil
}

func (s *Storage) getContract(address *evmInt256.Int) (*Contract, error) {
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

func (s *Storage) GetCode(address *evmInt256.Int) ([]byte, error) {
	contract, err := s.getContract(address)
	if err != nil {
		return nil, err
	}

	return contract.Code, err
}

func (s *Storage) GetCodeSize(address *evmInt256.Int) (*evmInt256.Int, error) {
	contract, err := s.getContract(address)
	if err != nil {
		return nil, err
	}

	return evmInt256.New(int64(contract.CodeSize)), err
}

func (s *Storage) GetCodeHash(address *evmInt256.Int) (*evmInt256.Int, error) {
	contract, err := s.getContract(address)
	if err != nil {
		return nil, err
	}

	return contract.CodeHash, err
}

func (s *Storage) GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error) {
	keyStr := block.AsStringKey()
	if hash, exists := s.readOnlyCache.BlockHash[keyStr]; exists {
		return hash, nil
	}

	hash, err := s.externalStorage.GetBlockHash(block)
	if err == nil {
		s.readOnlyCache.BlockHash[keyStr] = hash
	}

	return hash, err
}

func (s *Storage) NewContract(address *evmInt256.Int, code []byte) {
	s.ResultCache.NewContracts.Set(&Contract{
		Address:  address.Clone(),
		Code:     bytes.Clone(code),
		CodeHash: s.externalStorage.HashOfCode(code).Clone(),
		CodeSize: uint64(len(code)),
	})
}

func (s *Storage) CreateAddress(caller *evmInt256.Int, tx environment.Transaction) *evmInt256.Int {
	return s.externalStorage.CreateAddress(caller, tx)
}

func (s *Storage) CreateFixedAddress(caller *evmInt256.Int, salt *evmInt256.Int, code []byte, tx environment.Transaction) *evmInt256.Int {
	return s.externalStorage.CreateFixedAddress(caller, salt, code, tx)
}

func (s *Storage) GetExternalStorage() IExternalStorage {
	return s.externalStorage
}

func (s *Storage) ClearCache() {
	s.ResultCache = NewResultCache()
}
