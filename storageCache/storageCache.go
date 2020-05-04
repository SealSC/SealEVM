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
	Balance(address []byte) (*evmInt256.Int, error)
	GetCode(address []byte) ([]byte, error)
	GetCodeSize(address []byte) (*evmInt256.Int, error)
	GetCodeHash(address []byte) ([]byte, error)
	GetBlockHash(block uint64) ([]byte, error)

	load(k *evmInt256.Int) (*evmInt256.Int, error)
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
	ExternalStorage IExternalStorage
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
		ExternalStorage:     extStorage,
		ResultCallback: callback,
	}

	return s
}

func (s *StorageCache) SLoad(k *evmInt256.Int) (*evmInt256.Int, error ) {
	if s.ResultCache.OriginalData == nil || s.ResultCache.CachedData == nil || s.ExternalStorage == nil {
		return nil, evmErrors.StorageNotInitialized
	}

	i, exists := s.ResultCache.CachedData[k.String()]
	if exists {
		return i, nil
	}

	i, err := s.ExternalStorage.load(k)
	if err != nil {
		return nil, evmErrors.NoSuchDataInTheStorage(err)
	}

	cacheString := i.String()
	s.ResultCache.OriginalData[cacheString] = evmInt256.FromBigInt(i.Int)
	s.ResultCache.CachedData[cacheString] = i

	return i, nil
}

func (s *StorageCache) SStore(k *evmInt256.Int, v *evmInt256.Int)  {
	kString := k.String()
	s.ResultCache.CachedData[kString] = v
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
