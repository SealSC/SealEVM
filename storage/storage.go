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
	"github.com/SealSC/SealEVM/evmErrors"
	"github.com/SealSC/SealEVM/evmInt256"
)

type Storage struct {
	ResultCache     ResultCache
	ExternalStorage IExternalStorage
	readOnlyCache   readOnlyCache
}

func New(extStorage IExternalStorage) *Storage {
	s := &Storage{
		ResultCache: ResultCache{
			OriginalData: CacheUnderNamespace{},
			CachedData:   CacheUnderNamespace{},
			Balance:      BalanceCache{},
			Logs:         LogCache{},
			Destructs:    Cache{},
		},
		ExternalStorage: extStorage,
		readOnlyCache: readOnlyCache{
			Code:      CodeCache{},
			CodeSize:  Cache{},
			CodeHash:  Cache{},
			BlockHash: Cache{},
		},
	}

	return s
}

func (s *Storage) SLoad(n *evmInt256.Int, k *evmInt256.Int) (*evmInt256.Int, error ) {
	if s.ResultCache.OriginalData == nil || s.ResultCache.CachedData == nil || s.ExternalStorage == nil {
		return nil, evmErrors.StorageNotInitialized
	}

	nsStr := n.AsStringKey()
	keyStr := k.AsStringKey()

	var err error = nil
	i := s.ResultCache.CachedData.Get(nsStr, keyStr)
	if i == nil {
		i, err = s.ExternalStorage.Load(nsStr, keyStr)
		if err != nil {
			return nil, evmErrors.NoSuchDataInTheStorage(err)
		}

		s.ResultCache.OriginalData.Set(nsStr, keyStr, i)
		s.ResultCache.CachedData.Set(nsStr, keyStr, i)
	}

	return i, nil
}

func (s *Storage) SStore(n *evmInt256.Int, k *evmInt256.Int, v *evmInt256.Int)  {
	s.ResultCache.CachedData.Set(n.AsStringKey(), k.AsStringKey(), v)
}

func (s *Storage) BalanceModify(address *evmInt256.Int, value *evmInt256.Int, neg bool) {
	kString := address.AsStringKey()

	b, exist := s.ResultCache.Balance[kString]
	if !exist {
		b = &balance {
			Address: evmInt256.FromBigInt(address.Int),
			Balance: evmInt256.New(0),
		}

		s.ResultCache.Balance[kString] = b
	}

	if neg {
		b.Balance.Int.Sub(b.Balance.Int, value.Int)
	} else {
		b.Balance.Int.Add(b.Balance.Int, value.Int)
	}
}

func (s *Storage) Log(address *evmInt256.Int, topics [][]byte, data []byte, context environment.Context) {
	kString := address.AsStringKey()

	var theLog = Log{
		Topics:   topics,
		Data:    data,
		Context: context,
	}

	l := s.ResultCache.Logs[kString]
	s.ResultCache.Logs[kString] = append(l, theLog)

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

	return s.ExternalStorage.GetBalance(address)
}

func (s *Storage) GetCode(address *evmInt256.Int) ([]byte, error) {
	keyStr := address.AsStringKey()
	if b, exists := s.readOnlyCache.Code[keyStr]; exists {
		return b, nil
	}

	b, err := s.ExternalStorage.GetCode(address)
	if err == nil {
		s.readOnlyCache.Code[keyStr] = b
	}

	return b, err
}

func (s *Storage) GetCodeSize(address *evmInt256.Int) (*evmInt256.Int, error) {
	keyStr := address.AsStringKey()
	if size, exists := s.readOnlyCache.CodeSize[keyStr]; exists {
		return size, nil
	}

	size, err := s.ExternalStorage.GetCodeSize(address)
	if err == nil {
		s.readOnlyCache.CodeSize[keyStr] = size
	}

	return size, err
}

func (s *Storage) GetCodeHash(address *evmInt256.Int) (*evmInt256.Int, error) {
	keyStr := address.AsStringKey()
	if hash, exists := s.readOnlyCache.CodeHash[keyStr]; exists {
		return hash, nil
	}

	hash, err := s.ExternalStorage.GetCodeHash(address)
	if err == nil {
		s.readOnlyCache.CodeHash[keyStr] = hash
	}

	return hash, err
}

func (s *Storage) GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error) {
	keyStr := block.AsStringKey()
	if hash, exists := s.readOnlyCache.BlockHash[keyStr]; exists {
		return hash, nil
	}

	hash, err := s.ExternalStorage.GetBlockHash(block)
	if err == nil {
		s.readOnlyCache.BlockHash[keyStr] = hash
	}

	return hash, err
}
