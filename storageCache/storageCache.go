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

type IExternalStorage interface {
	Balance(address *evmInt256.Int) (*evmInt256.Int, error)
	GetCode(address *evmInt256.Int) ([]byte, error)
	GetCodeSize(address *evmInt256.Int) (*evmInt256.Int, error)
	GetCodeHash(address *evmInt256.Int) (*evmInt256.Int, error)
	GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error)

	CreateAddress(caller *evmInt256.Int) []byte
	CreateFixedAddress(caller *evmInt256.Int, salt *evmInt256.Int) []byte

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
	}

	return s
}

func (s *StorageCache) SLoad(n *evmInt256.Int, k *evmInt256.Int) (*evmInt256.Int, error ) {
	if s.ResultCache.OriginalData == nil || s.ResultCache.CachedData == nil || s.ExternalStorage == nil {
		return nil, evmErrors.StorageNotInitialized
	}

	cacheKey := n.String() + ":" +  k.String()
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
	cacheString := n.String() + "-" +  k.String()
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

//todo: cache all the below methods's result
func (s *StorageCache) Balance(address *evmInt256.Int) (*evmInt256.Int, error) {
	return s.ExternalStorage.Balance(address)
}

func (s *StorageCache) GetCode(address *evmInt256.Int) ([]byte, error) {
	//todo: handle precompiled-contract and any other special addresses code
	return s.ExternalStorage.GetCode(address)
}

func (s *StorageCache) GetCodeSize(address *evmInt256.Int) (*evmInt256.Int, error) {
	//todo: handle precompiled-contract and any other special addresses code size
	return s.ExternalStorage.GetCodeSize(address)
}

func (s *StorageCache) GetCodeHash(address *evmInt256.Int) (*evmInt256.Int, error) {
	//todo: handle precompiled-contract and any other special addresses code hash
	return s.ExternalStorage.GetCodeHash(address)
}

func (s *StorageCache) GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error) {
	return s.ExternalStorage.GetBlockHash(block)
}
