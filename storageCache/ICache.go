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
	SLoad(k *evmInt256.Int) (*evmInt256.Int, error)
	Balance(address []byte) (*evmInt256.Int, error)
	GetCode(address []byte) ([]byte, error)
	GetCodeSize(address []byte) (*evmInt256.Int, error)
	GetCodeHash(address []byte) ([]byte, error)
}

type ICache interface {
	SLoad(k *evmInt256.Int) (*evmInt256.Int, error)
	SStore(k *evmInt256.Int, i *evmInt256.Int)
	BalanceChange(address *evmInt256.Int, change *evmInt256.Int)
	Log(address *evmInt256.Int, topics [][]byte, data []byte, context environment.Context)
	Destruct(address *evmInt256.Int)
	ResultFeedback(err error)
}

type ResultCache struct {
	OriginalData    Cache
	CachedData      Cache
	Balance         BalanceCache
	Logs            LogCache
	Destructs       Cache
}

type storageCache struct {
	ResultCache

	extStorage      IExternalStorage
	callback        EVMResultCallback
}

func New(extStorage IExternalStorage, callback EVMResultCallback) ICache {
	s := &storageCache{
		extStorage: extStorage,
		callback: callback,
	}

	s.ResultCache.OriginalData = Cache{}
	s.ResultCache.CachedData = Cache{}
	s.ResultCache.Balance = BalanceCache{}
	s.ResultCache.Logs = LogCache{}
	s.ResultCache.Destructs = Cache{}

	return s
}

func (s *storageCache) SLoad(k *evmInt256.Int) (*evmInt256.Int, error ) {
	if s.OriginalData == nil || s.CachedData == nil || s.extStorage == nil {
		return nil, evmErrors.StorageNotInitialized
	}

	i, exists := s.CachedData[k.String()]
	if exists {
		return i, nil
	}

	i, err := s.extStorage.SLoad(k)
	if err != nil {
		return nil, evmErrors.NoSuchDataInTheStorage(err)
	}

	cacheString := i.String()
	s.OriginalData[cacheString] = evmInt256.FromBigInt(i.Int)
	s.CachedData[cacheString] = i

	return i, nil
}

func (s *storageCache) SStore(k *evmInt256.Int, v *evmInt256.Int)  {
	kString := k.String()
	s.CachedData[kString] = v
}

func (s *storageCache) BalanceChange(address *evmInt256.Int, change *evmInt256.Int) {
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

func (s *storageCache) Log(address *evmInt256.Int, topics [][]byte, data []byte, context environment.Context) {
	kString := address.String()

	var theLog = log {
		Topics:   topics,
		Data:    data,
		Context: context,
	}

	l := s.Logs[kString]
	s.Logs[kString] = append(l, theLog)

	return
}

func (s *storageCache) Destruct(address *evmInt256.Int) {
	s.ResultCache.Destructs[address.String()] = address
}

func (s *storageCache) ResultFeedback(err error) {
	s.callback(s.OriginalData, s.CachedData, err)
}
