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

type EVMResultCallback func(original Cache, final Cache, err error)

type IExternalStorage interface {
	balance(address *evmInt256.Int) (*evmInt256.Int, error)
	getCode(address *evmInt256.Int) ([]byte, error)
	getCodeSize(address *evmInt256.Int) (*evmInt256.Int, error)
	getCodeHash(address *evmInt256.Int) (*evmInt256.Int, error)
	getBlockHash(block *evmInt256.Int) (*evmInt256.Int, error)

	load(n *evmInt256.Int, k *evmInt256.Int) (*evmInt256.Int, error)
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
	ResultCallback  EVMResultCallback
	externalStorage IExternalStorage
}

func New(extStorage IExternalStorage, callback EVMResultCallback) *StorageCache {
	s := &StorageCache{
		ResultCache: ResultCache{
			OriginalData: Cache{},
			CachedData:   Cache{},
			Balance:      BalanceCache{},
			Logs:         LogCache{},
			Destructs:    Cache{},
		},
		externalStorage: extStorage,
		ResultCallback:  callback,
	}

	return s
}

func (s *StorageCache) SLoad(n *evmInt256.Int, k *evmInt256.Int) (*evmInt256.Int, error ) {
	if s.ResultCache.OriginalData == nil || s.ResultCache.CachedData == nil || s.externalStorage == nil {
		return nil, evmErrors.StorageNotInitialized
	}

	cacheKey := n.String() + ":" +  k.String()
	i, exists := s.ResultCache.CachedData[cacheKey]
	if exists {
		return i, nil
	}

	i, err := s.externalStorage.load(n, k)
	if err != nil {
		return nil, evmErrors.NoSuchDataInTheStorage(err)
	}

	s.ResultCache.OriginalData[cacheKey] = evmInt256.FromBigInt(i.Int)
	s.ResultCache.CachedData[cacheKey] = i

	return i, nil
}

func (s *StorageCache) SStore(n *evmInt256.Int, k *evmInt256.Int, v *evmInt256.Int)  {
	cacheString := n.String() + ":" +  k.String()
	s.ResultCache.CachedData[cacheString] = v
}

func (s *StorageCache) BalanceChange(address *evmInt256.Int, change *evmInt256.Int) {
	kString := address.String()

	b, exist := s.ResultCache.Balance[kString]
	if !exist {
		s.ResultCache.Balance[kString] = &balance {
			Address: address,
			Balance: change,
		}
	} else {
		b.Balance.Add(change)
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


//todo: cache all the below methods's result
func (s *StorageCache) Balance(address *evmInt256.Int) (*evmInt256.Int, error) {
	return s.externalStorage.balance(address)
}

func (s *StorageCache) GetCode(address *evmInt256.Int) ([]byte, error) {
	//todo: handle precompiled-contract and any other special addresses code
	return s.externalStorage.getCode(address)
}

func (s *StorageCache) GetCodeSize(address *evmInt256.Int) (*evmInt256.Int, error) {
	//todo: handle precompiled-contract and any other special addresses code size
	return s.externalStorage.getCodeSize(address)
}

func (s *StorageCache) GetCodeHash(address *evmInt256.Int) (*evmInt256.Int, error) {
	//todo: handle precompiled-contract and any other special addresses code hash
	return s.externalStorage.getCodeHash(address)
}

func (s *StorageCache) GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error) {
	return s.externalStorage.getBlockHash(block)
}
