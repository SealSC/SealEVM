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

package storageCache

import (
	"SealEVM/environment"
	"SealEVM/evmErrors"
	"SealEVM/evmInt256"
)

//todo: improve cache struct, because they has namespace-liked prefix.
type Cache map[string] *evmInt256.Int

type balance struct {
	Address *evmInt256.Int
	Balance *evmInt256.Int
}

type BalanceCache map[string] *balance

type log struct {
	Topics  [][]byte
	Data    []byte
	Context environment.Context
}
type LogCache map[string] []log

type IExternalStorage interface {
	GetBalance(address *evmInt256.Int) (*evmInt256.Int, error)
	GetCode(address *evmInt256.Int) ([]byte, error)
	GetCodeSize(address *evmInt256.Int) (*evmInt256.Int, error)
	GetCodeHash(address *evmInt256.Int) (*evmInt256.Int, error)
	GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error)

	CreateAddress(caller *evmInt256.Int) []byte
	CreateFixedAddress(caller *evmInt256.Int, salt *evmInt256.Int) []byte

	CanTransfer(from *evmInt256.Int, to *evmInt256.Int, amount *evmInt256.Int) bool

	Load(n *evmInt256.Int, k *evmInt256.Int) (*evmInt256.Int, error)
}

type ResultCache struct {
	OriginalData    Cache
	CachedData      Cache

	Balance         BalanceCache
	Logs            LogCache
	Destructs       Cache
}

type StorageCache struct {
	ResultCache     ResultCache
	ExternalStorage IExternalStorage
	readOnlyCache   readOnlyCache
}

type CodeCache map[string] []byte

type readOnlyCache struct {
	Code      CodeCache
	CodeSize  Cache
	CodeHash  Cache
	BlockHash Cache
}

func New(extStorage IExternalStorage) *StorageCache {
	s := &StorageCache{
		ResultCache: ResultCache{
			OriginalData: Cache{},
			CachedData:   Cache{},
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

func (s *StorageCache) SLoad(n *evmInt256.Int, k *evmInt256.Int) (*evmInt256.Int, error ) {
	if s.ResultCache.OriginalData == nil || s.ResultCache.CachedData == nil || s.ExternalStorage == nil {
		return nil, evmErrors.StorageNotInitialized
	}

	cacheKey := n.String() + "-" +  k.String()
	i, exists := s.ResultCache.CachedData[cacheKey]
	if exists {
		return i, nil
	}

	i, err := s.ExternalStorage.Load(n, k)
	if err != nil {
		return nil, evmErrors.NoSuchDataInTheStorage(err)
	}

	s.ResultCache.OriginalData[cacheKey] = evmInt256.FromBigInt(i.Int)
	s.ResultCache.CachedData[cacheKey] = i

	return i, nil
}

func (s *StorageCache) SStore(n *evmInt256.Int, k *evmInt256.Int, v *evmInt256.Int)  {
	cacheString := n.String() + "-" + k.String()
	s.ResultCache.CachedData[cacheString] = v
}

func (s *StorageCache) BalanceModify(address *evmInt256.Int, value *evmInt256.Int, neg bool) {
	kString := address.String()

	b, exist := s.ResultCache.Balance[kString]
	if !exist {
		s.ResultCache.Balance[kString] = &balance {
			Address: address,
			Balance: value,
		}
	}

	if neg {
		b.Balance.Sub(value)
	} else {
		b.Balance.Add(value)
	}
}

func (s *StorageCache) Log(address *evmInt256.Int, topics [][]byte, data []byte, context environment.Context) {
	kString := address.String()

	var theLog = log {
		Topics:   topics,
		Data:    data,
		Context: context,
	}

	l := s.ResultCache.Logs[kString]
	s.ResultCache.Logs[kString] = append(l, theLog)

	return
}

func (s *StorageCache) Destruct(address *evmInt256.Int) {
	s.ResultCache.Destructs[address.String()] = address
}

type commonGetterFunc func(*evmInt256.Int) (*evmInt256.Int, error)
func (s *StorageCache) commonGetter(key *evmInt256.Int, cache Cache, getterFunc commonGetterFunc) (*evmInt256.Int, error) {
	keyStr := key.String()
	if b, exists := cache[keyStr]; exists {
		return evmInt256.FromBigInt(b.Int), nil
	}

	b, err := getterFunc(key)
	if err == nil {
		cache[keyStr] = b
	}

	return b, err
}

func (s *StorageCache) Balance(address *evmInt256.Int) (*evmInt256.Int, error) {
	return s.ExternalStorage.GetBalance(address)
}

func (s *StorageCache) GetCode(address *evmInt256.Int) ([]byte, error) {
	keyStr := address.String()
	if b, exists := s.readOnlyCache.Code[keyStr]; exists {
		return b, nil
	}

	b, err := s.ExternalStorage.GetCode(address)
	if err == nil {
		s.readOnlyCache.Code[keyStr] = b
	}

	return b, err
}

func (s *StorageCache) GetCodeSize(address *evmInt256.Int) (*evmInt256.Int, error) {
	keyStr := address.String()
	if size, exists := s.readOnlyCache.CodeSize[keyStr]; exists {
		return size, nil
	}

	size, err := s.ExternalStorage.GetCodeSize(address)
	if err == nil {
		s.readOnlyCache.CodeSize[keyStr] = size
	}

	return size, err
}

func (s *StorageCache) GetCodeHash(address *evmInt256.Int) (*evmInt256.Int, error) {
	keyStr := address.String()
	if hash, exists := s.readOnlyCache.CodeHash[keyStr]; exists {
		return hash, nil
	}

	hash, err := s.ExternalStorage.GetCodeHash(address)
	if err == nil {
		s.readOnlyCache.CodeHash[keyStr] = hash
	}

	return hash, err
}

func (s *StorageCache) GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error) {
	keyStr := block.String()
	if hash, exists := s.readOnlyCache.BlockHash[keyStr]; exists {
		return hash, nil
	}

	hash, err := s.ExternalStorage.GetBlockHash(block)
	if err == nil {
		s.readOnlyCache.BlockHash[keyStr] = hash
	}

	return hash, err
}
